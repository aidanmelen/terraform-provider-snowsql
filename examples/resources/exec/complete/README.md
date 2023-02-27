# Complete Example

This examples shows how the `snowsql_exec` resource can manage [Snowflake](https://www.snowflake.com) objects using [SnowSQL](https://docs.snowflake.com/en/user-guide/snowsql.html).

## Usage

To run this example you need to execute:

```bash
$ terraform init
$ terraform plan
$ terraform apply
```

Note that this example may create resources which can cost money (Warehouse, Database Storage). Run `terraform destroy` when you don't need these resources.

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.13.0 |
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 2.1 |
| <a name="requirement_snowflake"></a> [snowflake](#requirement\_snowflake) | >= 0.56.5 |
| <a name="requirement_snowsql"></a> [snowsql](#requirement\_snowsql) | >= 1.3.2 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.4.3 |
| <a name="provider_snowflake"></a> [snowflake](#provider\_snowflake) | 0.56.5 |
| <a name="provider_snowsql"></a> [snowsql](#provider\_snowsql) | 1.3.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [random_password.password](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [random_pet.name](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [snowflake_database.database](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/database) | resource |
| [snowflake_role.role](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/role) | resource |
| [snowsql_exec.function](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/resources/exec) | resource |
| [snowsql_exec.role_grant_all](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/resources/exec) | resource |

## Inputs

No inputs.

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_show_function_results"></a> [show\_function\_results](#output\_show\_function\_results) | The SnowSQL query results from the read statements. |
| <a name="output_show_role_grant_all_results"></a> [show\_role\_grant\_all\_results](#output\_show\_role\_grant\_all\_results) | The SnowSQL query results from the read statements. |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
