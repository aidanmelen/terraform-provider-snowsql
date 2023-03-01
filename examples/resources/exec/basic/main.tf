resource "snowsql_exec" "role" {
  create {
    statements = "CREATE ROLE my_role"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role'"
  }

  # uncomment to update role in-place
  # update {
  #   statements = "ALTER ROLE my_role SET COMMENT = 'updated with terraform'"
  # }

  delete {
    statements = "DROP ROLE my_role"
  }
}
