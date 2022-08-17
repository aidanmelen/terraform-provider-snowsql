# Complete Example

Configuration in this directory demonstrates how to render and execute [DCL](https://www.geeksforgeeks.org/sql-ddl-dql-dml-dcl-tcl-commands/) SnowSQL commands from template files.

## Usage

To run this example you need to execute:

```bash
$ terraform init
$ terraform plan
$ terraform apply
```

Note that this example may create resources which can cost credits (Warehouse Usage, Database Storage).
Run `terraform destroy` when you don't need these resources.

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.13.0 |
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 2.1 |
| <a name="requirement_snowflake"></a> [snowflake](#requirement\_snowflake) | >= 0.33.4 |
| <a name="requirement_snowsql"></a> [snowsql](#requirement\_snowsql) | >= 0.4.3 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | >= 2.1 |
| <a name="provider_snowflake"></a> [snowflake](#provider\_snowflake) | >= 0.33.4 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_dcl"></a> [dcl](#module\_dcl) | ./modules/snowsql_exec_from_templates | n/a |

## Resources

| Name | Type |
|------|------|
| [random_pet.server](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [snowflake_database.database](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/database) | resource |
| [snowflake_database_grant.grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/database_grant) | resource |
| [snowflake_role.role](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role) | resource |
| [snowflake_role_grants.grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role_grants) | resource |
| [snowflake_schema.schema](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/schema) | resource |
| [snowflake_schema_grant.grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/schema_grant) | resource |
| [snowflake_table.table](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/table) | resource |
| [snowflake_user.user](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/user) | resource |
| [snowflake_warehouse.warehouse](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/warehouse) | resource |
| [snowflake_warehouse_grant.grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/warehouse_grant) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input\_name) | The name of the project. | `string` | `"COMPLETE_EXAMPLE"` | no |
| <a name="input_temporary_user_passworld"></a> [temporary\_user\_passworld](#input\_temporary\_user\_passworld) | The temporary password for the user. | `string` | `"ChangeMe2020!"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_snowsql_create_stmts"></a> [snowsql\_create\_stmts](#output\_snowsql\_create\_stmts) | The SnowSQL statements used during the resource lifecycle create. |
| <a name="output_snowsql_delete_stmts"></a> [snowsql\_delete\_stmts](#output\_snowsql\_delete\_stmts) | The SnowSQL statements used during the resource lifecycle delete. |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
