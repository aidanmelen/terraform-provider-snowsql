resource "snowsql_exec" "role" {
  name = local.name

  create {
    statements = "CREATE ROLE IF NOT EXISTS ${local.name};"
  }

  read {
    statements = "SHOW ROLES LIKE '${local.name}';"
  }

  delete {
    statements = "DROP ROLE IF EXISTS ${local.name};"
  }
}
