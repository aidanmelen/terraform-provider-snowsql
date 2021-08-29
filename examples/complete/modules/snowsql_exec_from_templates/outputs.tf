output "create_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle create."
  value       = snowsql_exec.stmts.create.0.statements
}

output "delete_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle delete."
  value       = snowsql_exec.stmts.delete.0.statements
}
