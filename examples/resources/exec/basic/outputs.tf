output "read_results" {
  description = "The SnowSQL query result from the read statements."
  value       = try(jsondecode(nonsensitive(snowsql_exec.role.read_results)), null)
}

output "create_results" {
  description = "The SnowSQL query result from the create statements."
  value       = nonsensitive(snowsql_exec.role.create_results)
}