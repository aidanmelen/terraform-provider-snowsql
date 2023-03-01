package snowsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/xid"
)

var numberOfStatementsDescription = "The number of SnowSQL statements to be executed. This can help reduce the risk of SQL injection attacks. Defaults to `null` indicating that there is no limit on the number of statements (`0` and `-1` also indicate no limit)."

var createLifecycleSchema = map[string]*schema.Schema{
	"statements": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: numberOfStatementsDescription,
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
		Description: numberOfStatementsDescription,
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
		Description: numberOfStatementsDescription,
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
				Optional:    true,
				ForceNew:    false,
				Default:     nil,
				Description: "The name of the resource. Defaults to random ID.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "" {
						errs = append(errs, fmt.Errorf("%q cannot be an empty string", key))
					}
					return
				},
			},
			"create": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				ForceNew:    true,
				Description: "Configuration block for create lifecycle statements.",
				Elem: &schema.Resource{
					Schema: createLifecycleSchema,
				},
			},
			"read": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ForceNew:    false,
				Description: "Configuration block for read lifecycle statements.",
				Elem: &schema.Resource{
					Schema: lifecycleSchema,
				},
			},
			"update": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ForceNew:    false,
				Description: "Configuration block for update lifecycle statements.",
				Elem: &schema.Resource{
					Schema: lifecycleSchema,
				},
			},
			"delete": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				ForceNew:    false,
				Description: "Configuration block for delete lifecycle statements.",
				Elem: &schema.Resource{
					Schema: lifecycleSchema,
				},
			},
			"read_results": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The encoded JSON list of query results from the read statements. This value is always marked as sensitive.",
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf("read_results", func(ctx context.Context, diff *schema.ResourceDiff, m interface{}) bool {
				return diff.HasChange("read.0.statements") || diff.HasChange("update.0.statements")
			}),
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceExecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	stmts := d.Get("create.0.statements").(string)
	numOfStmts := d.Get("create.0.number_of_statements").(int)

	db := m.(*sql.DB)
	err := snowflakeExecWithMultiStatement(ctx, db, "create", stmts, numOfStmts)
	if err != nil {
		return diag.FromErr(err)
	}

	name, ok := d.GetOk("name")
	if ok {
		d.SetId(name.(string))
	} else {
		id := xid.New().String()
		d.SetId(id)
	}

	diags = append(diags, resourceExecRead(ctx, d, m)...)
	return diags
}

func resourceExecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	stmts, ok := d.Get("read.0.statements").(string)
	numOfStmts := d.Get("read.0.number_of_statements").(int)

	if !ok || stmts == "" {
		d.Set("read", nil)
		d.Set("read_results", "null")
		return nil
	}

	db := m.(*sql.DB)
	results, err := snowflakeQueryWithMultiStatement(ctx, db, stmts, numOfStmts)
	if err != nil {
		d.Set("read", nil)
		d.Set("read_results", "null")
		return diag.FromErr(fmt.Errorf("failed to process the results from the query.\n\nStatements:\n\n  %s\n\nResults:\n\n  %v\n\n%s", stmts, results, err))
	}

	marshalledResults, _ := json.Marshal(results)
	if err != nil {
		d.Set("read", nil)
		d.Set("read_results", "null")
		return diag.FromErr(fmt.Errorf("failed to marshal query results to JSON.\n\nStatements:\n\n  %s\n\nResults:\n\n  %s\n\n%s", stmts, results, err))
	}

	log.Print("[DEBUG] marshalled query results: ", string(marshalledResults))

	d.Set("read_results", string(marshalledResults))

	return diags
}

func resourceExecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	stmts, ok := d.Get("update.0.statements").(string)
	numOfStmts := d.Get("update.0.number_of_statements").(int)

	if !ok || stmts == "" {
		d.Set("update", nil)
		return resourceExecRead(ctx, d, m)
	}

	if !d.HasChange("update.0.statements") {
		return resourceExecRead(ctx, d, m)
	}

	db := m.(*sql.DB)
	err := snowflakeExecWithMultiStatement(ctx, db, "update", stmts, numOfStmts)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceExecRead(ctx, d, m)
}

func resourceExecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	stmts := d.Get("delete.0.statements").(string)
	numOfStmts := d.Get("delete.0.number_of_statements").(int)

	db := m.(*sql.DB)
	err := snowflakeExecWithMultiStatement(ctx, db, "delete", stmts, numOfStmts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
