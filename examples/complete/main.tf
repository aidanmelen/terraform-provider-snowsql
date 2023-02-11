resource "random_pet" "name" {
  prefix    = upper(basename(path.cwd))
  separator = "_"
}

locals {
  name     = upper(random_pet.name.id)
}

###############################################################################
# Snowflake
###############################################################################
resource "snowflake_warehouse" "warehouse" {
  name = local.name
}

resource "snowflake_database" "database" {
  name = local.name
}

resource "snowflake_schema" "schema" {
  name     = local.name
  database = snowflake_database.database.name
}

resource "snowflake_table" "table" {
  name     = local.name
  database = snowflake_database.database.name
  schema   = snowflake_schema.schema.name

  column {
    name = "column1"
    type = "VARIANT"
  }
}

resource "snowflake_role" "role" {
  name = local.name
}

resource "snowflake_user" "user" {
  default_warehouse    = snowflake_warehouse.warehouse.name
  default_namespace    = join(".", [snowflake_database.database.name, snowflake_schema.schema.name])
  default_role         = snowflake_role.role.name
  must_change_password = true
  name                 = local.name
  password             = var.temporary_user_passworld
}

resource "snowflake_warehouse_grant" "grant" {
  warehouse_name = snowflake_warehouse.warehouse.name
  privilege      = "USAGE"
  roles          = [snowflake_role.role.name]
}

resource "snowflake_database_grant" "grant" {
  database_name = snowflake_database.database.name
  privilege     = "USAGE"
  roles         = [snowflake_role.role.name]
}

resource "snowflake_schema_grant" "grant" {
  database_name = snowflake_database.database.name
  privilege     = "USAGE"
  roles         = [snowflake_role.role.name]
  schema_name   = snowflake_schema.schema.name
}

###############################################################################
# SnowSQL
###############################################################################
resource "snowsql_exec" "grant_all" {
  name = local.name

  # grant all privileges on all (future) object when the resource is created
  create {
    statements = <<-EOT
      GRANT ALL PRIVILEGES ON ALL TABLES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON ALL VIEWS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON ALL FILE FORMATS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON ALL SEQUENCES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON ALL STREAMS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON ALL PROCEDURES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON FUTURE TABLES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON FUTURE VIEWS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON FUTURE FILE FORMATS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON FUTURE SEQUENCES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON FUTURE FUNCTIONS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON FUTURE STREAMS IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
      GRANT ALL PRIVILEGES ON FUTURE PROCEDURES IN DATABASE ${snowsql_exec.database.name} TO ROLE ${snowsql_exec.role.name};
    EOT
  }

  # revoke all grants when the resource is destroyed
  delete {
    statements = <<-EOT
      REVOKE ALL PRIVILEGES ON ALL TABLES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL VIEWS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL FILE FORMATS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL STREAMS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON ALL PROCEDURES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE TABLES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE VIEWS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE FILE FORMATS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE SEQUENCES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE FUNCTIONS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE STREAMS IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
			REVOKE ALL PRIVILEGES ON FUTURE PROCEDURES IN DATABASE ${snowsql_exec.database.name} FROM ROLE ${snowsql_exec.role.name};
    EOT
  }
}

resource "snowflake_role_grants" "grant" {
  role_name = snowflake_role.role.name
  users     = [snowflake_user.user.name]
}