# `snowsql_exec_from_templates` submodule

Helper submodule to read and render SnowSQL templates so they can be used to
manage the `create` and `delete` lifecycles controlled by the `snowsql_exec`
resource.

## Assumptions

This module assumes the existence of a `sql` directory containing the
SnowSQL `create` and `delete` templates. For examples, this is the default `sql`
file structure:

```text
workspace-root
├── main.tf
└── sql
    └── [NAME]
        ├── create.tpl
        └── delete.tpl
```

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.12.9 |
| <a name="requirement_snowsql"></a> [snowsql](#requirement\_snowsql) | >= 0.3.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_snowsql"></a> [snowsql](#provider\_snowsql) | >= 0.3.0 |
| <a name="provider_template"></a> [template](#provider\_template) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [snowsql_exec.stmts](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/resources/exec) | resource |
| [template_file.create](https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/file) | data source |
| [template_file.delete](https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/file) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input\_name) | The name of the directory containing the create and delete SnowSQL templates. | `string` | n/a | yes |
| <a name="input_sql_path"></a> [sql\_path](#input\_sql\_path) | The relative path to the `sql` directory containing the SnowSQL create and delete SnowSQL templates. Defaults to the `sql` directory in the root of your Terraform workspace. | `string` | `null` | no |
| <a name="input_template_vars"></a> [template\_vars](#input\_template\_vars) | The variables used to render the create and delete SnowSQL templates. This map will be passed to the [template\_file](https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/file#vars) data resource. | `map(any)` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_create_stmts"></a> [create\_stmts](#output\_create\_stmts) | The SnowSQL statements used during the resource lifecycle create. |
| <a name="output_delete_stmts"></a> [delete\_stmts](#output\_delete\_stmts) | The SnowSQL statements used during the resource lifecycle delete. |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
