resource "random_pet" "name" {
  prefix    = upper(basename(path.cwd))
  separator = "_"
}

locals {
  name = upper(random_pet.name.id)
}

###############################################################################
# Snowflake
###############################################################################
resource "snowflake_database" "database" {
  name = local.name
}

resource "snowflake_role" "role" {
  name = local.name
}

###############################################################################
# SnowSQL
###############################################################################
resource "snowsql_exec" "grant_all" {
  name = local.name

  # grant all privileges on all (future) objects when the resource is created
  create {
    statements = <<-EOT
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
  }

  # revoke all grants when the resource is destroyed
  delete {
    statements = <<-EOT
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
  }
}
