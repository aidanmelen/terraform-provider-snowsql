package snowsql

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUpdate(t *testing.T) {
	accName := fmt.Sprintf("terraform_provider_snowsql_acceptance_test_%s", acctest.RandStringFromCharSet(5, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: applyCreate(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowsql_exec.role", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.role", "create.0.statements", fmt.Sprintf("CREATE ROLE IF NOT EXISTS %s;", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "update.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "delete.0.statements", fmt.Sprintf("DROP ROLE IF EXISTS %s;", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "results", "example"),
				),
			},
			{
				Config: applyOptionalLifecycleBlocks(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowsql_exec.role", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.role", "create.0.statements", fmt.Sprintf("CREATE ROLE IF NOT EXISTS %s;", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read.%", fmt.Sprintf("SHOW ROLE ROLES LIKE %s;", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "update.0.statements", fmt.Sprintf("ALTER ROLE IF EXISTS %s SET COMMENT = 'updated with terraform';", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "delete.0.statements", fmt.Sprintf("DROP ROLE IF EXISTS %s;", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "results", "example"),
				),
			},
			{
				Config: destroyOptionalLifecycleBlocks(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowsql_exec.role", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.role", "create.0.statements", fmt.Sprintf("CREATE ROLE IF NOT EXISTS %s;", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "update.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "delete.0.statements", fmt.Sprintf("DROP ROLE IF EXISTS %s;", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "results", "example"),
				),
			},
		},
	})
}

func applyCreate(name string) string {
	s := `
	resource "snowsql_exec" "role" {
	  name = "%s"

		create {
			statements = "CREATE ROLE IF NOT EXISTS %s;"
		}

		delete {
			statements = "DROP ROLE IF EXISTS %s;"
		}
	}
	`
	return fmt.Sprintf(s, name, name, name)
}

func applyOptionalLifecycleBlocks(name string) string {
	s := `
	resource "snowsql_exec" "role" {
	  name = "%s"

		create {
			statements = "CREATE ROLE IF NOT EXISTS %s;"
		}

		create {
			statements = "SHOW ROLES LIKE %s;"
		}

		update {
		  	statements = "ALTER ROLE IF EXISTS %s SET COMMENT = 'updated with terraform';"
		}

		delete {
			statements = "DROP ROLE IF EXISTS %s;"
		}
	}
	`
	return fmt.Sprintf(s, name, name, name, name)
}

func destroyOptionalLifecycleBlocks(name string) string {
	s := `
	resource "snowsql_exec" "role" {
	  name = "%s"

		create {
			statements = "CREATE ROLE IF NOT EXISTS %s;"
		}

		delete {
			statements = "DROP ROLE IF EXISTS %s;"
		}
	}
	`
	return fmt.Sprintf(s, name, name, name)
}
