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
| terraform | >= 0.13.0 |
| random | >= 2.1 |
| snowflake | >= 0.25.18 |
| snowsql | >= 0.3.0 |

## Providers

| Name | Version |
|------|---------|
| random | >= 2.1 |
| snowflake | >= 0.25.18 |
| snowsql | >= 0.3.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| name | The name of the project. | `string` | `"SIMPLE_EXAMPLE"` | no |
| temporary\_user\_passworld | The temporary password for the user. | `string` | `"ChangeMe2020!"` | no |

## Outputs

| Name | Description |
|------|-------------|
| snowsql\_create\_stmts | The SnowSQL statements used during the resource lifecycle create. |
| snowsql\_delete\_stmts | The SnowSQL statements used during the resource lifecycle delete. |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
