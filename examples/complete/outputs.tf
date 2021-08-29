output "snowsql_create_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle create."
  value       = module.dcl.create_stmts
}

output "snowsql_delete_stmts" {
  description = "The SnowSQL statements used during the resource lifecycle delete."
  value       = module.dcl.delete_stmts
}
