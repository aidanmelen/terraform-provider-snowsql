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
	"github.com/snowflakedb/gosnowflake"
)

var numberOfStatementsDescription = "The number of SnowSQL statements to be executed. This can help reduce the risk of SQL injection attacks. Defaults to `null` indicating that there is no limit on the number of statements. `0` and `-1` also indicate no limit."

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
				Required:    true,
				ForceNew:    true,
				Description: "The name of the resource.",
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
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceExecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	createStmts := d.Get("create.0.statements").(string)
	numOfStmts := d.Get("create.0.number_of_statements").(int)

	// Execute the `create` statements
	db := m.(*sql.DB)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, createStmts)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to execute the create statements.\n\nStatements:\n\n  %s\n\n%s", createStmts, err))
	}

	name := d.Get("name").(string)
	d.SetId(name)

	return resourceExecRead(ctx, d, m)
}

func resourceExecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	readStmts, ok := d.Get("read.0.statements").(string)
	numOfStmts := d.Get("read.0.number_of_statements").(int)

	if !ok || readStmts == "" {
		d.Set("read", nil)
		d.Set("read_results", nil)
		return nil
	}

	// Query the `read` statements

	db := m.(*sql.DB)
	multiStmtCtx, err := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	if err != nil {
		d.Set("read", nil)
		d.Set("read_results", nil)
		return diag.FromErr(fmt.Errorf("failed to build multiple statement query.\n\nStatements:\n\n  %s\n\n%s", readStmts, err))
	}

	rows, err := db.QueryContext(multiStmtCtx, readStmts)

	if err != nil {
		d.Set("read", nil)
		d.Set("read_results", nil)
		return diag.FromErr(fmt.Errorf("failed to query the read statements.\n\nStatements:\n\n  %s\n\nResults:\n\n  %v\n\n%s", readStmts, rows, err))
	}
	defer rows.Close()

	// Process all the rows from all the queries and store the results in a list
	results := make([]map[string]interface{}, 0)
	processRows := func(rows *sql.Rows) error {
		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil {
				return err
			}
			values := make([]interface{}, len(columns))
			for i := range columns {
				values[i] = new(interface{})
			}
			err = rows.Scan(values...)
			if err != nil {
				return err
			}
			rowMap := make(map[string]interface{})
			for i, col := range columns {
				rowMap[col] = *values[i].(*interface{})
			}
			results = append(results, rowMap)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	}

	if err := processRows(rows); err != nil {
		d.Set("read", nil)
		d.Set("read_results", nil)
		return diag.FromErr(fmt.Errorf("failed to process the results from the first query.\n\nStatements:\n\n  %s\n\nResults:\n\n  %v\n\n%s", readStmts, results, err))
	}

	log.Print("[DEBUG] finished processing the first query result: ", results)

	for rows.NextResultSet() {
		if err := processRows(rows); err != nil {
			d.Set("read", nil)
			d.Set("read_results", nil)
			return diag.FromErr(fmt.Errorf("failed to process the query results.\n\nStatements:\n\n  %s\n\nResults:\n\n  %v\n\n%s", readStmts, results, err))
		}
	}

	log.Print("[DEBUG] finished processing the all query results: ", results)

	marshalledResults, _ := json.Marshal(results)
	if err != nil {
		d.Set("read", nil)
		d.Set("read_results", nil)
		return diag.FromErr(fmt.Errorf("failed to marshal query results to JSON.\n\nStatements:\n\n  %s\n\nResults:\n\n  %s\n\n%s", readStmts, results, err))
	}

	log.Print("[DEBUG] marshalled query results: ", string(marshalledResults))

	d.Set("read_results", string(marshalledResults))

	return diags
}

func resourceExecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	updateStmts, ok := d.Get("update.0.statements").(string)
	numOfStmts := d.Get("update.0.number_of_statements").(int)

	if !ok || updateStmts == "" {
		d.Set("update", nil)
		return resourceExecRead(ctx, d, m)
	}

	if !d.HasChange("update.0.statements") {
		return resourceExecRead(ctx, d, m)
	}

	// Execute the 'update' statements
	db := m.(*sql.DB)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, updateStmts)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to execute the update statements.\n\nStatements:\n\n  %s\n\n%s", updateStmts, err))
	}

	return resourceExecRead(ctx, d, m)
}

func resourceExecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	deleteStmts := d.Get("delete.0.statements").(string)
	numOfStmts := d.Get("delete.0.number_of_statements").(int)

	// Execute the 'delete' statements
	db := m.(*sql.DB)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, deleteStmts)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to execute the delete statements.\n\nStatements:\n\n  %s\n\n%s", deleteStmts, err))
	}

	d.SetId("")

	return diags
}
