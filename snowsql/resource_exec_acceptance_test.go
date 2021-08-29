package snowsql

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExec(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: execConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowsql_exec.database", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.role", "name", accName),
					resource.TestCheckResourceAttr("snowsql_exec.table_grants", "name", accName),
				),
			},
		},
	})
}

func execConfig(name string) string {
	s := `
	locals {
		name = "%s"
	}

	resource "snowsql_exec" "database" {
		name = local.name

		create {
			statements = "CREATE DATABASE IF NOT EXISTS ${local.name}"
		}

		delete {
			statements = "DROP DATABASE IF EXISTS ${local.name}"
		}

		delete_on_create = true
	}

	resource "snowsql_exec" "role" {
		name = local.name

		create {
			statements = "CREATE ROLE IF NOT EXISTS ${local.name}"
		}

		delete {
			statements = "DROP ROLE IF EXISTS ${local.name}"
		}

		delete_on_create = true
	}

	resource "snowsql_exec" "table_grants" {
		name = local.name

		create {
			statements = <<-EOT
			GRANT ALL PRIVILEGES ON ALL TABLES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON ALL VIEWS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON ALL FILE FORMATS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON ALL SEQUENCES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON ALL STREAMS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON ALL PROCEDURES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON FUTURE TABLES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON FUTURE VIEWS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON FUTURE FILE FORMATS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON FUTURE SEQUENCES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON FUTURE FUNCTIONS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON FUTURE STREAMS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			GRANT ALL PRIVILEGES ON FUTURE PROCEDURES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
			EOT
		}

		delete {
			statements =  <<-EOT
			REVOKE ALL PRIVILEGES ON ALL TABLES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL VIEWS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL FILE FORMATS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL STREAMS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL PROCEDURES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE TABLES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE VIEWS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE FILE FORMATS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE SEQUENCES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE FUNCTIONS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE STREAMS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE PROCEDURES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			EOT
		}

		delete_on_create = true
	}
	`
	return fmt.Sprintf(s, name)
}
