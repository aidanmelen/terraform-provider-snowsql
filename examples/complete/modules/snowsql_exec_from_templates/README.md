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
| terraform | >= 0.12.9 |
| snowsql | >= 0.1.0 |

## Providers

| Name | Version |
|------|---------|
| snowsql | >= 0.1.0 |
| template | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| name | The name of the directory containing the create and delete SnowSQL templates. | `string` | n/a | yes |
| sql\_path | The relative path to the `sql` directory containing the SnowSQL create and delete SnowSQL templates. Defaults to the `sql` directory in the root of your Terraform workspace. | `string` | `null` | no |
| template\_vars | The variables used to render the create and delete SnowSQL templates. This map will be passed to the [template\_file](https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/file#vars) data resource. | `map` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| create\_stmts | The SnowSQL statements used during the resource lifecycle create. |
| delete\_stmts | The SnowSQL statements used during the resource lifecycle delete. |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
