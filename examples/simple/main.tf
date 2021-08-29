resource "random_pet" "server" {
  prefix    = var.name
  separator = "_"
}

locals {
  name = upper(random_pet.server.id)
}

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

resource "snowsql_exec" "dcl" {
  name = local.name

  create {
    statements = <<-EOT
    GRANT ALL PRIVILEGES ON ALL TABLES IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
    GRANT ALL PRIVILEGES ON FUTURE TABLES IN DATABASE ${snowflake_database.database.name} TO ROLE ${snowflake_role.role.name};
    EOT
  }

  delete {
    statements = <<-EOT
    REVOKE ALL PRIVILEGES ON ALL TABLES IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
    REVOKE ALL PRIVILEGES ON FUTURE TABLES IN DATABASE ${snowflake_database.database.name} FROM ROLE ${snowflake_role.role.name};
    EOT
  }

  delete_on_create = true
}

resource "snowflake_role_grants" "grant" {
  role_name = snowflake_role.role.name
  users     = [snowflake_user.user.name]
}
