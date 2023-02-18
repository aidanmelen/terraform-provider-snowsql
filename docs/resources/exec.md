---
page_title: "snowsql_exec Resource - terraform-provider-snowsql"
subcategory: ""
description: |-
  
---

# snowsql_exec (Resource)

The `snowsql_exec` resource allows for custom Terraform CRUD management of [Snowflake](https://www.snowflake.com) objects using [SnowSQL](https://docs.snowflake.com/en/user-guide/snowsql.html).

## Examples

This basic example shows how to manage an arbitrary Snowflake object.

```terraform
resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role"
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role"
  }
}
```

-> **NOTE:** It is highly recommended to test all SnowSQL statements, especially create and delete statements, in a [Snowflake worksheet](https://docs.snowflake.com/en/user-guide/ui-worksheet) prior to executing them. This can help avoid any unexpected issues during the execution of these statements.

~> **NOTE:** It is important to ensure that any delete statements negate any corresponding create statements, to avoid any orphaned Snowflake objects. Failure to do so can result in clutter and potential issues within your Snowflake environment.

### Query Snowflake With Read Statements

This resource allows you to execute arbitrary SnowSQL queries and use the results in your infrastructure management.

```terraform
resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role'"
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role"
  }
}

output "show_role_results" {
  description = "The SnowSQL query result from the read statements."
  value       = jsondecode(nonsensitive(snowsql_exec.role.read_results))
}
```

This will output the JSON formatted results like this:

```bash
Outputs:

show_role_results = [
  {
    "assigned_to_users" = "0"
    "comment" = ""
    "created_on" = "2023-02-16T18:47:48.756-08:00"
    "granted_roles" = "0"
    "granted_to_roles" = "0"
    "is_current" = "N"
    "is_default" = "N"
    "is_inherited" = "N"
    "name" = "BASIC_CIVIL_HALIBUT"
    "owner" = "ACCOUNTADMIN"
  },
]
```

### Avoiding Replacement With Update Lifecycle

Execute the update statements as an in-place Terraform change by adding or modifying them after the initial Terraform apply. However, if there are any changes to the create statements, the resource will need to be replaced.

1. The create statements are run on the first apply:

```terraform
resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role'"
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role"
  }
}
```

2. Add the update statements to alter the role in-place.

```terraform
resource "snowsql_exec" "role" {
  name = local.name

  create {
    statements = "CREATE ROLE IF NOT EXISTS ${local.name}"
  }

  update {
    statements = "ALTER ROLE IF EXISTS ${local.name} SET COMMENT = 'updated with terraform'"
  }

  delete {
    statements = "DROP ROLE IF EXISTS ${local.name}"
  }
}
```

### Multi-Statements

This resource allows you to execute multiple SnowSQL statements separated by semicolons. This is particularly useful for managing multiple Snowflake objects within a single resource.

```terraform
resource "snowflake_database" "database" {
  name = local.name
}

resource "snowflake_role" "role" {
  name = local.name
}

resource "snowsql_exec" "role_grant_all" {
  name = local.name

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

  read {
    statements = <<-EOT
      SHOW GRANTS TO ROLE ${local.name};
      SHOW FUTURE GRANTS TO ROLE ${local.name};
    EOT
  }

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
```

-> **NOTE:** Please see the snowflakedb documentation for [Executing Multiple Statements in One Call](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#hdr-Executing_Multiple_Statements_in_One_Call) for more information.

~> **NOTE:** While it may be tempting to manage multiple statements with a single `snowsql_exec` resource, it can make it much harder to managed the resource in the future, as well as potentially causing issues if the resource needs to be rolled back. It is recommended to separate multiple statements into separate resources where possible.

### Multi-line Statements

Statements can also be formatted across multiple lines for better readability. This is particularly useful for executing complex SnowSQL statements.

```terraform
resource "snowsql_exec" "function" {
  name = local.name

  create {
    statements = <<-EOT
      CREATE OR REPLACE FUNCTION ${snowflake_database.database.name}.PUBLIC.JS_FACTORIAL(f FLOAT)
        RETURNS FLOAT
        LANGUAGE JAVASCRIPT
        STRICT
        AS '
        if (D <= 0) {
          return 1;
        } else {
          var result = 1;
          for (var i = 2; i <= D; i++) {
            result = result * i;
          }
          return result;
        }
        ';
    EOT
  }

  read {
    statements = <<-EOT
      SHOW USER FUNCTIONS LIKE 'JS_FACTORIAL' 
        IN DATABASE ${snowflake_database.database.name};
    EOT
  }

  delete {
    statements = <<-EOT
      DROP FUNCTION IF EXISTS 
        ${snowflake_database.database.name}.PUBLIC.JS_FACTORIAL(FLOAT);
    EOT
  }
}
```

## Argument Reference

* `name` - (Required, Forces new resource) The name of the resource.
* `create` - (Required, Forces new resource) Configuration block for create lifecycle statements. (see [below for nested schema](#nestedblock-lifecycle))
* `read` - (Optional) Configuration block for read lifecycle statements. (see [below for nested schema](#nestedblock-lifecycle))
* `update` - (Optional) Configuration block for in-place update lifecycle statements. (see [below for nested schema](###nestedblock-lifecycle))
* `delete` - (Required) Configuration block for delete lifecycle statements. (see [below for nested schema](#nestedblock-lifecycle))

### Nested Blocks - `create`, `read`, `update`, and `delete`

The nested blocks all have the same arguments.

- `statements` (Required) A string containing one or many SnowSQL statements separated by semicolons.
- `number_of_statements` (Optional) The number of SnowSQL statements. This can help reduce the risk of SQL injection attacks.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `read_results` (String) The encoded JSON list of query results from the read statements. This value is always marked as sensitive.

## Import

Import is supported using the following syntax:

```shell
terraform import snowsql_exec.name name
```