resource "snowsql_exec" "role" {
  name = local.name

  create {
    statements = "CREATE ROLE IF NOT EXISTS ${local.name};"
  }

  read {
    statements = "SHOW ROLES LIKE '${local.name}';"
  }

  update {
    statements = "ALTER ROLE IF EXISTS ${local.name} SET COMMENT = 'updated with terraform';"
  }

  delete {
    statements = "DROP ROLE IF EXISTS ${local.name};"
  }
}
