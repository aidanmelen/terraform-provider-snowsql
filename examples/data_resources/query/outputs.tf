output "show_role_results" {
  description = "The SnowSQL query result from the read statements."
  value       = jsondecode(nonsensitive(data.snowsql_query.accountadmin_role.results))
}