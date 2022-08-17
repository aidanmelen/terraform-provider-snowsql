
resource "random_pet" "server" {
  prefix    = "GRANTS"
  separator = "_"
}

locals {
  name_grants = upper(random_pet.server.id)
  update_grant_all_database = join(
    "",
    flatten(
      [
        for id,db in snowflake_table_grant.update_db_future : formatlist("GRANT UPDATE, INSERT, DELETE, TRUNCATE ON ALL TABLES IN DATABASE ${db.database_name} TO ROLE %s;", db.roles)
      ]
    )
  )
  update_grant_all_schema = join(
    "",
    flatten(
      [
        for id,db in snowflake_table_grant.update_schema_future : formatlist("GRANT UPDATE, INSERT, DELETE, TRUNCATE ON ALL TABLES IN SCHEMA ${db.database_name}.${db.schema_name} TO ROLE %s;", db.roles)
      ]
    )
  )
}

resource "snowsql_exec" "db_update" {
  name = "${local.name_grants}_ALL_DATABASES"

  create {
    statements = <<-EOT
    ${local.update_grant_all_database}
    EOT
  }
  delete {
    statements = <<-EOT
    REVOKE SELECT PRIVILEGES ON ALL MATERIALIZED VIEWS IN DATABASE TEST TO ROLE TEST;
    EOT
  }
}

resource "snowsql_exec" "sch_update" {
  name = "${local.name_grants}_ALL_SCHEMAS"

  create {
    statements = <<-EOT
    ${local.update_grant_all_schema}
    EOT
  }

  delete {
    statements = <<-EOT
    REVOKE SELECT PRIVILEGES ON ALL MATERIALIZED VIEWS IN DATABASE TEST TO ROLE TEST;
    EOT
  }
}
