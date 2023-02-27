output "show_role_results" {
  description = "The SnowSQL query result from the read statements."
  value       = jsondecode(nonsensitive(snowsql_exec.role.read_results))
}
