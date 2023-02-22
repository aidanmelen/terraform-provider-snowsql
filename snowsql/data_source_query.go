package snowsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/xid"
)

func dataSourceQuery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQueryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     nil,
				Description: "The name of the data resource. Defaults to random ID.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "" {
						errs = append(errs, fmt.Errorf("%q cannot be an empty string", key))
					}
					return
				},
			},
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
				Description: "The number of SnowSQL statements to be executed. This can help reduce the risk of SQL injection attacks. Defaults to `null` indicating that there is no limit on the number of statements (`0` and `-1` also indicate no limit).",
			},
			"results": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The encoded JSON list of query results from the query statements. This value is always marked as sensitive.",
			},
		},
	}
}

func dataSourceQueryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name, nameOk := d.GetOk("name")
	stmts := d.Get("statements").(string)
	numOfStmts := d.Get("number_of_statements").(int)

	db := m.(*sql.DB)
	results, err := snowflakeQueryWithMultiStatement(ctx, db, stmts, numOfStmts)
	if err != nil {
		d.Set("results", nil)
		return diag.FromErr(fmt.Errorf("failed to process the results from the query.\n\nStatements:\n\n  %s\n\nResults:\n\n  %v\n\n%s", stmts, results, err))
	}

	marshalledResults, _ := json.Marshal(results)
	if err != nil {
		d.Set("results", nil)
		return diag.FromErr(fmt.Errorf("failed to marshal query results to JSON.\n\nStatements:\n\n  %s\n\nResults:\n\n  %s\n\n%s", stmts, results, err))
	}

	log.Print("[DEBUG] marshalled query results: ", string(marshalledResults))

	d.Set("results", string(marshalledResults))

	if nameOk {
		d.SetId(name.(string))
	} else {
		id := xid.New().String()
		d.SetId(id)
	}

	return diags
}
