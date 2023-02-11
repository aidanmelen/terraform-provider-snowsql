# Simple Example

Configuration in this directory demonstrates how to execute a [DCL](https://www.geeksforgeeks.org/sql-ddl-dql-dml-dcl-tcl-commands/) with the `snowsql_exec` resource.

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
| <a name="requirement_snowflake"></a> [snowflake](#requirement\_snowflake) | >= 0.33.4 |
| <a name="requirement_snowsql"></a> [snowsql](#requirement\_snowsql) | >= 1.2.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.4.3 |
| <a name="provider_snowsql"></a> [snowsql](#provider\_snowsql) | 1.2.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [random_pet.name](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/pet) | resource |
| [snowsql_exec.role](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/resources/exec) | resource |

## Inputs

No inputs.

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_snowsql_create_stmts"></a> [snowsql\_create\_stmts](#output\_snowsql\_create\_stmts) | The SnowSQL statements used during the resource lifecycle create. |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
