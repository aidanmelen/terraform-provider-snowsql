resource "snowflake_database" "database" {
  name = local.name
}

resource "snowflake_role" "role" {
  name = local.name
}

resource "snowsql_exec" "role_grant_all" {
  name = local.name

  create {
    statements           = <<-EOT
      GRANT ALL PRIVILEGES ON ALL TABLES IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON ALL VIEWS IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON ALL FILE FORMATS IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON ALL SEQUENCES IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON ALL STREAMS IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON ALL PROCEDURES IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON FUTURE TABLES IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON FUTURE VIEWS IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON FUTURE FILE FORMATS IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON FUTURE SEQUENCES IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON FUTURE FUNCTIONS IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON FUTURE STREAMS IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
      GRANT ALL PRIVILEGES ON FUTURE PROCEDURES IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
    EOT
    number_of_statements = 14
  }

  read {
    statements           = <<-EOT
      SHOW GRANTS TO ROLE ${local.name};
      SHOW FUTURE GRANTS TO ROLE ${local.name};
    EOT
    number_of_statements = 2
  }

  delete {
    statements           = <<-EOT
      REVOKE ALL PRIVILEGES ON ALL TABLES IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON ALL VIEWS IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON ALL FILE FORMATS IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON ALL STREAMS IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON ALL PROCEDURES IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON FUTURE TABLES IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON FUTURE VIEWS IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON FUTURE FILE FORMATS IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON FUTURE SEQUENCES IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON FUTURE FUNCTIONS IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON FUTURE STREAMS IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
      REVOKE ALL PRIVILEGES ON FUTURE PROCEDURES IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
    EOT
    number_of_statements = 14
  }
}
