resource "snowsql_exec" "role" {
  name = local.name

  # create snowflake object(s) during the resource creation
  create {
    statements = "CREATE ROLE IF NOT EXISTS ${local.name};"
  }

  read {
    statements = "SHOW ROLES LIKE '${local.name}';"
  }

  # uncomment to alter the snowflake object(s) during the in-place resource change
  update {
    statements = "ALTER ROLE IF EXISTS ${local.name} SET COMMENT = 'updated with terraform';"
  }

  # drop the snowflake object(s) during the resource destruction
  delete {
    statements = "DROP ROLE IF EXISTS ${local.name};"
  }
}
