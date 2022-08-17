output "snowsql_sch_update" {
  description = "The SnowSQL statements used during the resource lifecycle create databases."
  value       = module.db_update.create_stmts
}

output "snowsql_db_update" {
  description = "The SnowSQL statements used during the resource lifecycle create schema."
  value       = module.sch_update.create_stmts
}
