output "show_roles_like_accountadmin" {
  description = "Show the ACCOUNTADMIN role."
  value       = jsondecode(nonsensitive(data.snowsql_query.show_roles_like_accountadmin.results))
}

output "show_database_like_sample" {
  description = "Show all Snowflake sample databases."
  value       = jsondecode(nonsensitive(data.snowsql_query.show_database_like_sample.results))
}

output "count_snowflake_sample_data_tables" {
  description = "Count some tables from snowflake_sample_data"
  value       = jsondecode(nonsensitive(data.snowsql_query.count_snowflake_sample_data_tables.results))
}

output "select_snowflake_sample_data_tpch_sf1_lineitem" {
  description = "select complex query from snowflake_sample_data.tpch_sf1.lineitem."
  value       = jsondecode(nonsensitive(data.snowsql_query.select_snowflake_sample_data_tpch_sf1_lineitem.results))
}

