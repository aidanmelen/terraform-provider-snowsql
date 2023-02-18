resource "snowsql_exec" "role" {
  name = local.name

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role;"
  }

  delete {
    statements = "DROP ROLE IF EXISTS ${local.name};"
  }
}

