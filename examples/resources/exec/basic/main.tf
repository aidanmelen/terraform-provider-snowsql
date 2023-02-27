resource "snowsql_exec" "role" {
  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role'"
  }

  # uncomment to update role in-place
  # update {
  #   statements = "ALTER ROLE IF EXISTS my_role SET COMMENT = 'updated with terraform'"
  # }

  delete {
    statements = "DROP ROLE IF EXISTS my_role"
  }
}
