output "results" {
  description = "The SnowSQL query result from the read statements."
  value       = try(jsondecode(nonsensitive(snowsql_exec.role.results)), null)
}