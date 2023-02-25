---
page_title: "Migrating `snowsql_exec` resources to snowflake provider."
description: |-
  This guide explains how to migrate `snowsql_exec` resources to the Snowflake provider in Terraform.
---

This guide explains how to migrate `snowsql_exec` resources to the Snowflake provider in Terraform, specifically for granting ownership of a Snowflake user to a role.

## Step 1: Create a Snowflake User

Start by creating a new Snowflake user with the `snowflake_user` resource:

```terraform
provider "snowflake" {}

resource "snowflake_user" "user" {
  name = "MY_USER"
}
```

-> **Note** The ownership of the Snowflake user will be determined by [discretionary access control](https://docs.snowflake.com/en/user-guide/security-access-control-overview#access-control-framework).

## Step 2: Grant Role Ownership using Snowsql

For sake of example, let's assume that the [Snowflake provider does not yet support a resource for granting user ownership to a role](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/956) at the time of creation. We will need to use a custom `snowsql_exec` resource to grant user ownership to the `USERADMIN` role:

```terraform
provider "snowflake" {}
provider "snowsql" {}

resource "snowflake_user" "user" {
  name = "MY_USER"
}

resource "snowsql_exec" "user_ownership_grant" {
    name = "snowsql_user_ownership_grant"

    create {
        statements = "GRANT OWNERSHIP ON USER ${snowflake_user.user.name} TO ROLE SYSADMIN COPY CURRENT GRANTS;"
    }

    delete {
        statements = "GRANT OWNERSHIP ON USER ${snowflake_user.user.name} TO ROLE ACCOUNTADMIN COPY CURRENT GRANTS;"
    }
}
```

## Step 3: Use Snowflake Provider for User Ownership Grant

Now that the [Snowflake provider supports a `snowflake_user_ownership_grant` resource](https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/969), we can simplify the Terraform configuration by replacing the `snowsql_exec` resource with the new `snowflake_user_ownership_grant` resource:


```terraform
provider "snowflake" {}

resource "snowflake_user" "user" {
  name = "MY_USER"
}

resource "snowflake_user_ownership_grant" "grant" {
	on_user_name                  = snowflake_user.user.name
	to_role_name                  = "USERADMIN"
	current_grants                = "COPY"
  revert_ownership_to_role_name = "ACCOUNTADMIN"
}
```

Once you've updated the terraform, remove the old `snowsql_exec` resource from the Terraform state:

```console
$ terraform state rm snowsql_exec.user_ownership_grant
Removed snowsql_exec.user_ownership_grant
Successfully removed 1 resource instance(s).
```

Next, import the new Snowflake provider resource using the following command:

```console
$ terraform import snowflake_user_ownership_grant.grant 'MY_USER|USERADMIN|COPY'
```

The output should confirm that the resource was imported successfully:

```
snowflake_user_ownership_grant.grant: Importing from ID "MY_USER|USERADMIN|COPY"...
snowflake_user_ownership_grant.grant: Import prepared!
  Prepared snowflake_user_ownership_grant for import
snowflake_user_ownership_grant.grant: Refreshing state... [id=MY_USER|USERADMIN|COPY]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```