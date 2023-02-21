resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role;"
  }

  read {
    statements = "SHOW ROLES LIKE 'SYSADMIN';\nSHOW ROLES LIKE 'ACCOUNTADMIN';"
    # statements = <<-EOF
    #   SHOW ROLES LIKE 'ACCOUNTADMIN';
    #   SHOW ROLES LIKE 'my_role';
    # EOF
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role"
  }
}