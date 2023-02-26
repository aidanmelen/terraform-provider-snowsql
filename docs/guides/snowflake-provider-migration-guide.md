---
page_title: "Migrating `snowsql_exec` resources to snowflake provider."
description: |-
  This guide explains how to migrate `snowsql_exec` resources to the Snowflake provider in Terraform.
---

This guide provides step-by-step instructions for migrating SnowSQL resources to the Snowflake provider in Terraform, specifically for migrating a Snowflake user.

## Step 1: Create a Snowflake User with the SnowSQL Provider

- Begin by creating a new Snowflake user using the `snowsql_exec` resource.
- Add the following code to your Terraform file, which specifies the SnowSQL provider and creates a `snowsql_exec` resource:

```terraform
terraform {
  required_version = ">= 0.13.0"

  required_providers {
    snowsql = {
      source  = "aidanmelen/snowsql"
      version = ">= 1.3.1"
    }
  }
}

provider "snowsql" {}

resource "snowsql_exec" "example_user" {
    name = "my_user"

    create {
        statements = <<-EOT
          CREATE USER MY_USER 
            WITH COMMENT = 'created with terraform-provider-snowsql';
        EOT
    }

    read {
      statements = "SHOW USERS LIKE 'MY_USER';"
    }

    delete {
        statements = "DROP USER MY_USER;"
    }
}

output "user_comment" {
  description = "The Snowflake user comment."
  value       = lookup(jsondecode(nonsensitive(snowsql_exec.example_user.read_results))[0], "comment", null)
}
```

- Run `terraform apply` to create the user and retrieve user information from Snowflake.
- The output should display the user comment you specified in the `snowsql_exec` resource.

## Step 2: Migrate User to the Snowflake Provider

Update the Terraform file to use the Snowflake provider and create a `snowflake_user` resource instead of the `snowsql_exec` resource:

```terraform
terraform {
  required_version = ">= 0.13.0"

  required_providers {
      snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = ">= 0.56.5"
    }
  }
}

provider "snowflake" {}

resource "snowflake_user" "example" {
  name = "MY_USER"
  comment = "created with terraform-provider-snowsql"
}

output "user_comment" {
  description = "The Snowflake user comment."
  value       = snowflake_user.example.comment
}
```

- Remove the old `snowsql_exec` resource from the Terraform state by running `terraform state rm snowsql_exec.example_user`.
- Import the new `snowflake_user` resource into the Terraform state by running `terraform import snowflake_user.example 'MY_USER'`.
- Run `terraform apply` to output the original user comment from the `snowflake_user` resource.

The migration is now complete.