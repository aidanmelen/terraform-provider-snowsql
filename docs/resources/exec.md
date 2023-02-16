---
page_title: "snowsql_exec Resource - terraform-provider-snowsql"
subcategory: ""
description: |-
  
---

# snowsql_exec (Resource)



## Example Basic Usage

```terraform
resource "snowsql_exec" "role" {
  name = local.name

  create {
    statements = "CREATE ROLE IF NOT EXISTS ${local.name};"
  }

  read {
    statements = "SHOW ROLES LIKE '${local.name}';"
  }

  delete {
    statements = "DROP ROLE IF EXISTS ${local.name};"
  }
}
```

## Example Multi-Statement Usage

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

## Avoiding Replacement

Any changes to the `create` statements will cause a replacement change. 
Adding or changing the `update` statements will result in an in-place change with the execution of the `update` statement.

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

    **NOTE** the `create` statements are only executed on creation or when the statements change.

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
- `read_results` (String, Sensitive) The List of read query results.

<a id="nestedblock--create"></a>
### Nested Schema for `create`

Required:

- `statements` (String) A string containing one or many SnowSQL statements separated by semicolons. it's worth noting that splitting queries in this way is not always reliable since some SQL statements (e.g., CREATE FUNCTION) can contain semicolons within the statement itself.

Optional:

- `number_of_statements` (Number) A string containing one or many SnowSQL statements separated by semicolons. it's worth noting that splitting queries in this way is not always reliable since some SQL statements (e.g., CREATE FUNCTION) can contain semicolons within the statement itself.


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