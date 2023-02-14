output "read_results" {
  description = "The SnowSQL query result from the read statements."
  value       = try(jsondecode(nonsensitive(snowsql_exec.role.read_results)), null)
}

output "read_results_raw" {
  description = "The SnowSQL query result from the read statements."
  value       = nonsensitive(snowsql_exec.role.read_results)
}