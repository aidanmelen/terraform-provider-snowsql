package snowsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/snowflakedb/gosnowflake"
)

var numberOfStatementsDescription = "A string containing one or many SnowSQL statements separated by semicolons. it's worth noting that splitting queries in this way is not always reliable since some SQL statements (e.g., CREATE FUNCTION) can contain semicolons within the statement itself."

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
		ForceNew:    true,
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
			"read_results": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The List of read query results.",
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
	numOfStmts, ok := l[0].(map[string]interface{})["number_of_statements"].(int)

	if !ok {
		numOfStmts = strings.Count(multiStmt, ";")
	}

	return multiStmt, numOfStmts
}

func resourceExecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Execute the `create` statements
	db := m.(*sql.DB)
	name := d.Get("name").(string)
	multiStmt, numOfStmts := parseLifecycleSchemaData("create", d)
	multiStmtCtx, _ := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	_, err := db.ExecContext(multiStmtCtx, multiStmt)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceExecRead(ctx, d, m)
}

func resourceExecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	readStmts, ok := d.Get("read.0.statements").(string)
	if !ok {
		d.Set("read_results", nil)
		return nil
	}

	// Execute the `read` query statements
	db := m.(*sql.DB)
	sdb := sqlx.NewDb(db, "snowflake").Unsafe()
	rows, err := sdb.Queryx(readStmts)

	if err != nil {
		fmt.Println("Error running query:", err)
		d.Set("read_results", nil)
		return diags
	}

	// Loop through the query result rows
	var queryResult []map[string]interface{}
	for rows.Next() {
		rowData := make(map[string]interface{})
		err := rows.MapScan(rowData)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			d.Set("read_results", nil)
			return diag.FromErr(err)
		}
		// Append the row data to the query result array
		queryResult = append(queryResult, rowData)
	}

	if err := rows.Close(); err != nil {
		d.Set("read_results", nil)
		return diag.FromErr(err)
	}

	log.Print("[DEBUG] raw query result: ", queryResult)

	// Marshal the query read_results to JSON
	marshalledResult, err := json.Marshal(queryResult)
	if err != nil {
		d.Set("read_results", nil)
		return diag.FromErr(err)
	}

	log.Print("[DEBUG] marshalled query result: ", string(marshalledResult))

	d.Set("read_results", string(marshalledResult))

	return diags
}

func resourceExecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	_, ok := d.GetOk("update.0.statements")
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
		return diag.FromErr(err)
	}

	return resourceExecRead(ctx, d, m)
}

func resourceExecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Execute the 'delete' statements
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
