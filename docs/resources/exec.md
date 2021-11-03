---
page_title: "snowsql_exec Resource - terraform-provider-snowsql"
subcategory: ""
description: |-
  A resource for executing arbitrary SnowSQL statements.
---

# Resource: snowsql_exec

A resource for executing arbitrary SnowSQL statements.

## Example Usage

This demonstrates how to execute a [DCL](https://www.geeksforgeeks.org/sql-ddl-dql-dml-dcl-tcl-commands/) with the `snowsql_exec` resource.

-> We use the snowflake provider to manage the snowflake resources and we use snowsql_exec to execute verbose table grants that are not supported by the snowflake provider, as of 08/29/2021.

```hcl
terraform {
  required_version = ">= 0.13.0"

  required_providers {
    snowflake = {
      source  = "chanzuckerberg/snowflake"
      version = ">= v0.25.18"
    }
    snowsql = {
      source  = "aidanmelen/snowsql"
      version = ">= 0.2.0"
    }
    random = ">= 2.1"
  }
}

provider "snowflake" {}
provider "snowsql" {}

resource "random_pet" "server" {
  prefix    = "EXAMPLE"
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

output "snowsql_create_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle create."
  value       = snowsql_exec.dcl.create.0.statements
}

output "snowsql_delete_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle delete."
  value       = snowsql_exec.dcl.delete.0.statements
}
```

## Argument Reference

* `name` - (Required) Specifies the identifier for the SnowSQL commands.
* `create` - (Required) Specifies the SnowSQL create lifecycle. See [Lifecycle item](#lifecycle-item) below for details.
* `delete` - (Required) Specifies the SnowSQL delete lifecycle. See [Lifecycle item](#lifecycle-item) below for details.
* `delete_on_create` - (Optional) Execute delete statements before create statements during the create lifecycle. See [Delete On Create Example](#delete-on-create-example) below for details.

### Lifecycle item

-> We recommend testing the SnowSQL statements in a [Snowflake worksheet](https://docs.snowflake.com/en/user-guide/ui-worksheet.html) prior to automating with Terraform.

!> Failure to ensure that delete statements negate create statements may lead to unexpected consquences.

Each lifecyle item contains `statements` and `number_of_statements`.

- `statements` - (Required) A string containing one or many SnowSQL statements separated by semicolons. See [Understanding Multiple Statement Failures](#understanding-multiple-statement-failures) below for details.
- `number_of_statements` - (Optional) Specifies the number of SnowSQL statements. Defaults to `-1` which will dynamically count the number semicolons in SnowSQL statements.

## Understanding Multiple Statement Failures

If any of the SQL statements fail to compile or execute, execution is aborted. Any previous statements that ran before are unaffected.

For example, if the statements below are run as one multi-statement query, the multi-statement query fails on the third statement, and an exception is thrown.

```sql
CREATE OR REPLACE TABLE test(n int);
INSERT INTO TEST VALUES (1), (2);
INSERT INTO TEST VALUES ('not_an_integer');  -- execution fails here
INSERT INTO TEST VALUES (3);
```

If you then query the contents of the table named "test", the values 1 and 2 would be present. See [gosnowflake](https://godoc.org/github.com/snowflakedb/gosnowflake#hdr-Executing_Multiple_Statements_in_One_Call) for more details.

## Delete On Create Example

-> A *zombie resource* is a resource that is created in terraform state but never dies. It requires manual intervention to remove from state.

`delete_on_create` ensures that both delete and create statements compile and execute before the resource is applied to state. The following example will fail on apply and prevent the creation of a zombie resource.

```hcl
resource "snowsql_exec" "dcl" {
  name = "dcl"

  create {
    statements = <<-EOT
    CREATE USER IF NOT EXISTS example
    COMMENT 'hopefully the delete statements negate me.';
    EOT
  }

  delete {
    statements = "DROPPER USER IF EXISTS example; -- this will fail"
  }

  delete_on_create = true
}
```

This is logically equivalent to

```hcl
resource "snowsql_exec" "dcl" {
  name = "dcl"

  create {
    statements = <<-EOT
    DROPPER USER IF EXISTS EXAMPLE; -- this will fail
    CREATE USER IF NOT EXISTS example
    COMMENT 'hopefully the delete statements negate me.';
    EOT
  }

  delete {
    statements = "DROPPER USER IF EXISTS example; -- this will fail"
  }
}
```


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The identifier for the SnowSQL commands.
* `name` - The identifier for the SnowSQL commands.
* `create` - The SnowSQL create lifecycle.
* `delete` - The SnowSQL delete lifecycle.
