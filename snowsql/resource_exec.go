package snowsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
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
		Default:     nil,
		Computed:    true,
		Description: "Specifies the number of SnowSQL statements. If not provided, the default value is the count of semicolons in SnowSQL statements.",
	},
}

var lifecycleSchema = map[string]*schema.Schema{
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
		Default:     nil,
		Computed:    true,
		Description: "Specifies the number of SnowSQL statements. If not provided, the default value is the count of semicolons in SnowSQL statements.",
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
				Description: "Specifies the identifier for the SnowSQL resource.",
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
			"read": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ForceNew:    false,
				Description: "Specifies the SnowSQL read lifecycle.",
				Elem: &schema.Resource{
					Schema: lifecycleSchema,
				},
			},
			// TODO: https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/customizing-differences
			"update": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ForceNew:    false,
				Description: "Specifies the SnowSQL update lifecycle.",
				Elem: &schema.Resource{
					Schema: lifecycleSchema,
				},
			},
			"delete": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				ForceNew:    false,
				Description: "Specifies the SnowSQL delete lifecycle.",
				Elem: &schema.Resource{
					Schema: lifecycleSchema,
				},
			},
			"results": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The read query results.",
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
	numOfStmts, ok := l[0].(map[string]interface{})["number_of_statements"].(int)

	if !ok {
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

	resourceExecRead(ctx, d, m)

	d.SetId(name)

	return diags
}

func resourceExecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	readStmts, ok := d.Get("read.0.statements").(string)
	if !ok {
		d.Set("results", nil)
		return nil
	}

	db := m.(*sql.DB)

	sdb := sqlx.NewDb(db, "snowflake").Unsafe()
	rows, err := sdb.Queryx(readStmts)

	if err != nil {
		fmt.Println("Error running query:", err)
		return diags
	}

	var resultData []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		err := rows.MapScan(row)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return diag.FromErr(err)
		}
		resultData = append(resultData, row)
	}

	if err := rows.Err(); err != nil {
		return diag.FromErr(err)
	}

	log.Print("[DEBUG] results ", resultData)
	jsonResultData, err := json.Marshal(resultData)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Print("[DEBUG] results json ", string(jsonResultData))
	if err := d.Set("results", string(jsonResultData)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceExecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, ok := d.GetOk("update.0.statements")
	if !ok {
		d.Set("update", nil)
		return nil
	}

	if !d.HasChange("update.0.statements") {
		return nil
	}

	db := m.(*sql.DB)
	multiStmt, numOfStmts := parseLifecycleSchemaData("update", d)

	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, multiStmt)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceExecRead(ctx, d, m)

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
