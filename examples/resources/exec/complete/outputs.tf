output "show_role_grant_all" {
  description = "The SnowSQL query results from the read statements."
  value       = try(jsondecode(nonsensitive(snowsql_exec.role_grant_all.read_results)), null)
}

output "show_function" {
  description = "The SnowSQL query results from the read statements."
  value       = try(jsondecode(nonsensitive(snowsql_exec.function.read_results)), null)
}