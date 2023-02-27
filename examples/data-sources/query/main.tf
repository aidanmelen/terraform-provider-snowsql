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
data "snowsql_query" "select_snowflake_sample_data_tpch_sf1_lineitem" {
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
