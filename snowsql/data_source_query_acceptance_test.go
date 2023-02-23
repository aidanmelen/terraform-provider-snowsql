package snowsql

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceQuery_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQueryConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowsql_query.example", "result"),
				),
			},
			{
				Config: testAccDataSourceQueryConfigWithOptionals(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowsql_query.example", "result"),
				),
			},
		},
	})
}

func testAccDataSourceQueryConfig() string {
	return `
		data "snowsql_query" "example" {
			statements = "SELECT current_user()"
		}
	`
}

func testAccDataSourceQueryConfigWithOptionals() string {
	return `
		data "snowsql_query" "example" {
			name       			 = "example"
			statements 		     = "SELECT current_user()"
			number_of_statements = 1
		}
	`
}
