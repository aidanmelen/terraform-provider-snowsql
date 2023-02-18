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

var numberOfStatementsDescription = "The number of SnowSQL statements. This can help reduce the risk of SQL injection attacks. Defaults to `null`."

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

func parseLifecycleSchemaData(lifecycle string, d *schema.ResourceData) (string, int) {
	l := d.Get(lifecycle).([]interface{})
	multiStmt := l[0].(map[string]interface{})["statements"].(string)
	numOfStmts := l[0].(map[string]interface{})["number_of_statements"].(int)
	return multiStmt, numOfStmts
}

func resourceExecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	createStmts := d.Get("create.0.statements")

	// Execute the `create` statements
	db := m.(*sql.DB)
	multiStmt, numOfStmts := parseLifecycleSchemaData("create", d)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, multiStmt)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to execute the create statements.\n\nStatements:\n\n  %s\n\n%s", createStmts, err))
	}

	name := d.Get("name").(string)
	d.SetId(name)

	return resourceExecRead(ctx, d, m)
}

func resourceExecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, ok := d.Get("read.0.statements").(string)
	if !ok {
		d.Set("read_results", nil)
		return nil
	}

	// Execute the `create` statements
	db := m.(*sql.DB)
	multiStmt, numOfStmts := parseLifecycleSchemaData("read", d)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	rows, err := db.QueryContext(multiStmtCtx, multiStmt)
	if err != nil {
		d.Set("read", nil)
		d.Set("read_results", nil)
		return diag.FromErr(err)
	}
	defer rows.Close()

	// Get the names of the columns in the result set
	columnNames, err := rows.Columns()
	if err != nil {
		d.Set("read", nil)
		d.Set("read_results", nil)
		return diag.FromErr(err)
	}

	// Loop through the query result rows
	var queryResult []map[string]interface{}
	for rows.Next() {
		// Create a []interface{} slice with the same length as the number of columns
		columnValues := make([]interface{}, len(columnNames))
		// Create a []interface{} slice to hold the values for this row
		rowValues := make([]interface{}, len(columnNames))
		for i := range rowValues {
			rowValues[i] = &columnValues[i]
		}
		// Scan the row values into the columnValues slice
		if err := rows.Scan(rowValues...); err != nil {
			d.Set("read", nil)
			d.Set("read_results", nil)
			return diag.FromErr(err)
		}
		// Create a map of column name to value for this row
		row := make(map[string]interface{})
		for i, name := range columnNames {
			row[name] = columnValues[i]
		}
		queryResult = append(queryResult, row)
	}

	log.Print("[DEBUG] raw query result: ", queryResult)

	// Marshal the query read_results to JSON
	marshalledResult, err := json.Marshal(queryResult)
	if err != nil {
		d.Set("read", nil)
		d.Set("read_results", nil)
		return diag.FromErr(err)
	}

	log.Print("[DEBUG] marshalled query result: ", string(marshalledResult))

	d.Set("read_results", string(marshalledResult))

	return diags
}

func resourceExecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	updateStmts, ok := d.GetOk("update.0.statements")
	if !ok {
		d.Set("update", nil)
		return resourceExecRead(ctx, d, m)
	}

	if !d.HasChange("update.0.statements") {
		return resourceExecRead(ctx, d, m)
	}

	// Execute the 'update' statements
	db := m.(*sql.DB)
	multiStmt, numOfStmts := parseLifecycleSchemaData("update", d)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, multiStmt)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to execute the update statements.\n\nStatements:\n\n  %s\n\n%s", updateStmts, err))
	}

	return resourceExecRead(ctx, d, m)
}

func resourceExecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	deleteStmts := d.Get("delete.0.statements")

	// Execute the 'delete' statements
	db := m.(*sql.DB)
	multiStmt, numOfStmts := parseLifecycleSchemaData("delete", d)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, multiStmt)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to execute the delete statements.\n\nStatements:\n\n  %s\n\n%s", deleteStmts, err))
	}

	d.SetId("")

	return diags
}
