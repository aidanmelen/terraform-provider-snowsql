output "read_results" {
  description = "The SnowSQL query results from the read statements."
  value       = try(jsondecode(nonsensitive(snowsql_exec.role.read_results)), null)
}