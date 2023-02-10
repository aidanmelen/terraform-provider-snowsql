package snowsql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowflakedb/gosnowflake"
)

var createLifecycleSchema = map[string]*schema.Schema{
	"statements": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "A string containing one or many SnowSQL statements separated by semicolons.",
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			v := val.(string)
			if v == "" {
				errs = append(errs, fmt.Errorf("%q cannot be an empty string", key))
			}
			return
		},
	},
	"number_of_statements": {
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    true,
		Default:     -1,
		Description: "Specifies the number of SnowSQL statements. Defaults to `-1` which will dynamically count the number semicolons in SnowSQL statements. Go [here](https://godoc.org/github.com/snowflakedb/gosnowflake#hdr-Executing_Multiple_Statements_in_One_Call) to learn more about preventing SQL injection attacks.",
	},
}

var updateLifecycleSchema = map[string]*schema.Schema{
	"statements": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    false,
		Description: "A string containing one or many SnowSQL statements separated by semicolons.",
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			v := val.(string)
			if v == "" {
				errs = append(errs, fmt.Errorf("%q cannot be an empty string", key))
			}
			return
		},
	},
	"number_of_statements": {
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    false,
		Default:     -1,
		Description: "Specifies the number of SnowSQL statements. Defaults to `-1` which will dynamically count the number semicolons in SnowSQL statements. Go [here](https://godoc.org/github.com/snowflakedb/gosnowflake#hdr-Executing_Multiple_Statements_in_One_Call) to learn more about preventing SQL injection attacks.",
	},
}

var deleteLifecycleSchema = map[string]*schema.Schema{
	"statements": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    false,
		Description: "A string containing one or many SnowSQL statements separated by semicolons.",
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			v := val.(string)
			if v == "" {
				errs = append(errs, fmt.Errorf("%q cannot be an empty string", key))
			}
			return
		},
	},
	"number_of_statements": {
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    false,
		Default:     -1,
		Description: "Specifies the number of SnowSQL statements. Defaults to `-1` which will dynamically count the number semicolons in SnowSQL statements. Go [here](https://godoc.org/github.com/snowflakedb/gosnowflake#hdr-Executing_Multiple_Statements_in_One_Call) to learn more about preventing SQL injection attacks.",
	},
}

func resourceExec() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExecCreate,
		ReadContext:   resourceExecRead,
		UpdateContext: resourceExecUpdate,
		DeleteContext: resourceExecDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the identifier for the SnowSQL commands.",
			},
			"create": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				ForceNew:    true,
				Description: "Specifies the SnowSQL create lifecycle.",
				Elem: &schema.Resource{
					Schema: createLifecycleSchema,
				},
			},
			"update": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ForceNew:    false,
				Description: "Specifies the SnowSQL update lifecycle.",
				Elem: &schema.Resource{
					Schema: updateLifecycleSchema,
				},
			},
			"delete": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				ForceNew:    false,
				Description: "Specifies the SnowSQL delete lifecycle.",
				Elem: &schema.Resource{
					Schema: deleteLifecycleSchema,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func parseLifecycleSchemaData(lifecycle string, d *schema.ResourceData) (string, int) {
	l := d.Get(lifecycle).([]interface{})
	multiStmt := l[0].(map[string]interface{})["statements"].(string)
	numOfStmts := l[0].(map[string]interface{})["number_of_statements"].(int)

	if numOfStmts == -1 {
		numOfStmts = strings.Count(multiStmt, ";")
	}

	return multiStmt, numOfStmts
}

func resourceExecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	db := m.(*sql.DB)
	name := d.Get("name").(string)

	multiStmt, numOfStmts := parseLifecycleSchemaData("create", d)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, multiStmt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return diags
}

// resourceExecRead is not implemented and stubbed out because it is required.
func resourceExecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceExecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// if update is set, execute the given lifecycle commands
	if _, ok := d.GetOk("update"); ok {
		db := m.(*sql.DB)
		multiStmt, numOfStmts := parseLifecycleSchemaData("update", d)

		multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
		_, err := db.ExecContext(multiStmtCtx, multiStmt)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceExecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	db := m.(*sql.DB)
	multiStmt, numOfStmts := parseLifecycleSchemaData("delete", d)

	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, multiStmt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
