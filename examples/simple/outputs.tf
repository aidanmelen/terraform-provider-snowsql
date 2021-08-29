output "snowsql_create_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle create."
  value       = snowsql_exec.dcl.create.0.statements
}

output "snowsql_delete_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle delete."
  value       = snowsql_exec.dcl.delete.0.statements
}
