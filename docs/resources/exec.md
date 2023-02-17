---
page_title: "snowsql_exec Resource - terraform-provider-snowsql"
subcategory: ""
description: |-
  
---

# snowsql_exec (Resource)

The `snowsql_exec` resource allows for the management of the `create`, `read`, `update`, and `delete` lifecycles for [Snowflake](https://www.snowflake.com) objects using [SnowSQL](https://docs.snowflake.com/en/user-guide/snowsql.html).

## Examples

This example shows how to `create` and `delete` a Snowflake role with the `snowsql_exec` resource.

```terraform
resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role;"
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role;"
  }
}
```

-> **NOTE:** It is highly recommended to test all SnowSQL statements, especially `create` and `delete` statements, in a [Snowflake worksheet](https://docs.snowflake.com/en/user-guide/ui-worksheet) prior to executing them. This can help avoid any unexpected issues during the execution of these statements.

!> **WARNING:** It is important to ensure that any `delete` statements negate any corresponding `create` statements, to avoid any orphaned Snowflake objects. Failure to do so can result in clutter and potential issues within your Snowflake environment.

### Query Snowflake With Read Statements

The `snowsql_exec` resource allows you to execute arbitrary Snowflake SQL queries from Terraform, and use the results in your infrastructure management. When using the `read` statements in the resource, the result(s) of the SQL query/queries will be available in the `read_results` attribute as a sensitive raw JSON string. 

To output the query results in a non-sensitive format to the console, you can use the `nonsensitive` function to mark the `read_results` value as non-sensitive before decoding it with the `jsondecode` function.

```terraform
resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role;"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role';"
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role;"
  }
}

output "read_results" {
  description = "The SnowSQL query result from the read statements."
  value       = jsondecode(nonsensitive(snowsql_exec.role.read_results))
}
```

This will output the JSON formatted results like this:

```bash
Outputs:

read_results = [
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

This example demonstrates the behavior of the `snowsql_exec` resource when using the optional `update` statements to modify an existing object. If the `update` statements are added or changed after the initial terraform apply, Terraform will perform an in-place change by executing the `update` statement(s). On the other hand, any changes to the `create` statements will cause a replacement of the resource.

1. The `create` statements are run on the first apply:

```terraform
resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role;"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role';"
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role;"
  }
}
```

2. Add the `update` statements to alter the role in-place.

```terraform
resource "snowsql_exec" "role" {
  name = "my_role"

  create {
    statements = "CREATE ROLE IF NOT EXISTS my_role;"
  }

  read {
    statements = "SHOW ROLES LIKE 'my_role';"
  }

  update {
    statements = "ALTER ROLE IF EXISTS my_role SET COMMENT = 'updated with terraform';"
  }

  delete {
    statements = "DROP ROLE IF EXISTS my_role;"
  }
}
```

### Multi-Statements

The `snowsql_exec` resource allows you to execute multiple SnowSQL statements separated by semicolons. This is particularly useful for managing multiple Snowflake objects within a single `snowsql_exec` resource.

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

All `snowsql_exec` resource statements can be formatted across multiple lines for better readability. This is particularly useful for executing complex SnowSQL statements.

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
    statements = "SHOW USER FUNCTIONS LIKE 'JS_FACTORIAL' IN DATABASE ${snowflake_database.database.name};"
  }

  delete {
    statements = "DROP FUNCTION IF EXISTS ${snowflake_database.database.name}.PUBLIC.JS_FACTORIAL(FLOAT);"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `create` (Block List, Min: 1, Max: 1) Specifies the SnowSQL create lifecycle. (see [below for nested schema](#nestedblock--create))
- `delete` (Block List, Min: 1, Max: 1) Specifies the SnowSQL delete lifecycle. (see [below for nested schema](#nestedblock--delete))
- `name` (String) Specifies the identifier for the SnowSQL resource.

### Optional

- `read` (Block List, Max: 1) Specifies the SnowSQL read lifecycle. (see [below for nested schema](#nestedblock--read))
- `update` (Block List, Max: 1) Specifies the SnowSQL update lifecycle. (see [below for nested schema](#nestedblock--update))

### Read-Only

- `id` (String) The ID of this resource.
- `read_results` (String, Sensitive) The List of query results from the read statements.

<a id="nestedblock--create"></a>
### Nested Schema for `create`

Required:

- `statements` (String) A string containing one or many SnowSQL statements separated by semicolons.

Optional:

- `number_of_statements` (Number) A string containing one or many SnowSQL statements separated by semicolons.


<a id="nestedblock--delete"></a>
### Nested Schema for `delete`

Required:

- `statements` (String) A string containing one or many SnowSQL statements separated by semicolons.

Optional:

- `number_of_statements` (Number) Specifies the number of SnowSQL statements. If not provided, the default value is the count of semicolons in SnowSQL statements.


<a id="nestedblock--read"></a>
### Nested Schema for `read`

Required:

- `statements` (String) A string containing one or many SnowSQL statements separated by semicolons.

Optional:

- `number_of_statements` (Number) Specifies the number of SnowSQL statements. If not provided, the default value is the count of semicolons in SnowSQL statements.


<a id="nestedblock--update"></a>
### Nested Schema for `update`

Required:

- `statements` (String) A string containing one or many SnowSQL statements separated by semicolons.

Optional:

- `number_of_statements` (Number) Specifies the number of SnowSQL statements. If not provided, the default value is the count of semicolons in SnowSQL statements.

## Import

Import is supported using the following syntax:

```shell
terraform import snowsql_exec.name name
```