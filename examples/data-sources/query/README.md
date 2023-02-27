# Query Example

This examples shows how the `snowsql_query` data source can retrieve information about [Snowflake](https://www.snowflake.com) objects using [SnowSQL](https://docs.snowflake.com/en/user-guide/snowsql.html).

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
| <a name="provider_snowsql"></a> [snowsql](#provider\_snowsql) | >= 1.3.2 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [snowsql_query.count_snowflake_sample_data_tables](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/data-sources/query) | data source |
| [snowsql_query.select_current_user](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/data-sources/query) | data source |
| [snowsql_query.select_snowflake_sample_data_tpch_sf1_lineitem](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/data-sources/query) | data source |
| [snowsql_query.show_database_like_sample](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/data-sources/query) | data source |

## Inputs

No inputs.

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_count_snowflake_sample_data_tables"></a> [count\_snowflake\_sample\_data\_tables](#output\_count\_snowflake\_sample\_data\_tables) | Count some tables from snowflake\_sample\_data |
| <a name="output_select_current_user"></a> [select\_current\_user](#output\_select\_current\_user) | Select the current user. |
| <a name="output_select_snowflake_sample_data_tpch_sf1_lineitem"></a> [select\_snowflake\_sample\_data\_tpch\_sf1\_lineitem](#output\_select\_snowflake\_sample\_data\_tpch\_sf1\_lineitem) | select complex query from snowflake\_sample\_data.tpch\_sf1.lineitem. |
| <a name="output_show_database_like_sample"></a> [show\_database\_like\_sample](#output\_show\_database\_like\_sample) | Show all Snowflake sample databases. |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
