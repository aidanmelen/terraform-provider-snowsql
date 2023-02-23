---
page_title: "snowsql_query Data Source - terraform-provider-snowsql"
subcategory: ""
description: "The `snowsql_query` data resource."
---

# snowsql_query (Data Source)

The `snowsql_query` data resource allow you to retrieve information from [Snowflake](https://www.snowflake.com) objects using [SnowSQL](https://docs.snowflake.com/en/user-guide/snowsql.html).

## Examples

This example shows how to query arbitrary Snowflake objects.

```terraform
data "snowsql_query" "select_current_user" {
  statements = "SELECT current_user()"
}

data "snowsql_query" "show_database_like_sample" {
  statements = "SHOW DATABASES LIKE '%sample%'"
}

# multi-statement queries
data "snowsql_query" "count_snowflake_sample_data_tables" {
  statements = <<-EOT
    select 'customer' AS table_name, count(*) as count from snowflake_sample_data.tpch_sf1.customer;
    select 'lineitem' AS table_name, count(*) as count from snowflake_sample_data.tpch_sf1.lineitem;
    select 'nation' AS table_name, count(*) as count from snowflake_sample_data.tpch_sf1.nation;
  EOT
}

# multi-line statement query
data  "snowsql_query" "select_snowflake_sample_data_tpch_sf1_lineitem" {
  statements = <<-EOT
    // https://docs.snowflake.com/en/user-guide/sample-data-tpch#functional-query-definition
    
    use schema snowflake_sample_data.tpch_sf1;

    select
        l_returnflag,
        l_linestatus,
        sum(l_quantity) as sum_qty,
        sum(l_extendedprice) as sum_base_price,
        sum(l_extendedprice * (1-l_discount)) as sum_disc_price,
        sum(l_extendedprice * (1-l_discount) * (1+l_tax)) as sum_charge,
        avg(l_quantity) as avg_qty,
        avg(l_extendedprice) as avg_price,
        avg(l_discount) as avg_disc,
        count(*) as count_order
    from
        lineitem
    where
        l_shipdate <= dateadd(day, -90, to_date('1998-12-01'))
    group by
        l_returnflag,
        l_linestatus
    order by
        l_returnflag,
        l_linestatus
    limit 10;
  EOT
}
```

-> **NOTE:** It is highly recommended to test all SnowSQL query statements in a [Snowflake worksheet](https://docs.snowflake.com/en/user-guide/ui-worksheet) prior to executing them. This can help avoid any unexpected issues during the execution of these statements.

-> **NOTE:** The query statements are executed and the resulting row(s) are processed in the same way as the [`snowsql_exec` read](https://registry.terraform.io/providers/aidanmelen/snowsql/latest/docs/resources/exec#query-snowflake-with-read-statements) statements.

## Argument Reference

* `name` - (Required) The name of the resource.
- `statements` (Required) A string containing one or many SnowSQL statements separated by semicolons.
- `number_of_statements` (Optional) The number of SnowSQL statements to be executed. This can help reduce the risk of SQL injection attacks. Defaults to `null` indicating that there is no limit on the number of statements (`0` and `-1` also indicate no limit).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `results` (String) The encoded JSON list of query results from the query statements. This value is always marked as sensitive.