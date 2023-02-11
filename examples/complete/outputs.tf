output "snowsql_create_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle create."
  value       = snowsql_exec.grant_all.create.0.statements
}
