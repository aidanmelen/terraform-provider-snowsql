output "show_role_grant_all_results" {
  description = "The SnowSQL query results from the read statements."
  value       = jsondecode(nonsensitive(snowsql_exec.role_grant_all.read_results))
}

output "show_function_results" {
  description = "The SnowSQL query results from the read statements."
  value       = jsondecode(nonsensitive(snowsql_exec.function.read_results))
}
