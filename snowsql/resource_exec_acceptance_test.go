package snowsql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLifecycles(t *testing.T) {
	accName := fmt.Sprintf("TERRAFORM_PROVIDER_SNOWSQL_ACCEPTANCE_TEST_%s", strings.ToUpper(acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)))

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: applyCreate(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowsql_exec.role", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.role", "create.0.statements", fmt.Sprintf("CREATE ROLE IF NOT EXISTS %s", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "update.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "delete.0.statements", fmt.Sprintf("DROP ROLE IF EXISTS %s", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read_results", "null"),
				),
			},
			{
				Config: applyOptionalLifecycleBlocks(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowsql_exec.role", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.role", "create.0.statements", fmt.Sprintf("CREATE ROLE IF NOT EXISTS %s", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read.0.statements", fmt.Sprintf("SHOW ROLES LIKE '%s'", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "update.0.statements", fmt.Sprintf("ALTER ROLE IF EXISTS %s SET COMMENT = 'updated with terraform'", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "delete.0.statements", fmt.Sprintf("DROP ROLE IF EXISTS %s", accName)),
					resource.TestCheckResourceAttrSet("snowsql_exec.role", "read_results"),
				),
			},
			{
				Config: destroyOptionalLifecycleBlocks(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowsql_exec.role", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.role", "create.0.statements", fmt.Sprintf("CREATE ROLE IF NOT EXISTS %s", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "update.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "delete.0.statements", fmt.Sprintf("DROP ROLE IF EXISTS %s", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read_results", "null"),
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
			statements = "CREATE ROLE IF NOT EXISTS %s"
		}

		delete {
			statements = "DROP ROLE IF EXISTS %s"
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
			statements = "CREATE ROLE IF NOT EXISTS %s"
		}

		read {
			statements = "SHOW ROLES LIKE '%s'"
		}

		update {
		  	statements = "ALTER ROLE IF EXISTS %s SET COMMENT = 'updated with terraform'"
		}

		delete {
			statements = "DROP ROLE IF EXISTS %s"
		}
	}
	`
	return fmt.Sprintf(s, name, name, name, name, name)
}

func destroyOptionalLifecycleBlocks(name string) string {
	s := `
	resource "snowsql_exec" "role" {
	  name = "%s"

		create {
			statements = "CREATE ROLE IF NOT EXISTS %s"
		}

		delete {
			statements = "DROP ROLE IF EXISTS %s"
		}
	}
	`
	return fmt.Sprintf(s, name, name, name)
}

func TestAccRead(t *testing.T) {
	accName := fmt.Sprintf("TERRAFORM_PROVIDER_SNOWSQL_ACCEPTANCE_TEST_%s", strings.ToUpper(acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)))

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: readLifecycleBlockWithMultipleStatements(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowsql_exec.role", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.role", "create.0.statements", fmt.Sprintf("CREATE ROLE IF NOT EXISTS %s", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "read.0.statements", fmt.Sprintf("SHOW ROLES LIKE '%s';\nSHOW ROLES \n\tLIKE 'SYSADMIN';\n", accName)),
					resource.TestCheckResourceAttr("snowsql_exec.role", "update.%", "0"),
					resource.TestCheckResourceAttr("snowsql_exec.role", "delete.0.statements", fmt.Sprintf("DROP ROLE IF EXISTS %s", accName)),
					resource.TestCheckResourceAttrSet("snowsql_exec.role", "read_results"),
				),
			},
		},
	})
}

func readLifecycleBlockWithMultipleStatements(name string) string {
	s := `
	resource "snowsql_exec" "role" {
	  name = "%s"

		create {
			statements = "CREATE ROLE IF NOT EXISTS %s"
		}

		read {
			statements = <<-EOT
				SHOW ROLES LIKE '%s';
				SHOW ROLES
					LIKE 'SYSADMIN';
			EOT
		}

		delete {
			statements = "DROP ROLE IF EXISTS %s"
		}
	}
	`
	return fmt.Sprintf(s, name, name, name, name)
}
