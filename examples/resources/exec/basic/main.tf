resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role'"
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role"
  }
}