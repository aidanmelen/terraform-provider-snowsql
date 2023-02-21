resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role'"
    number_of_statements = 1
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role"
  }
}