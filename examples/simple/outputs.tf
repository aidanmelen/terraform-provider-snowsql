output "snowsql_create_stmts" {
  description = "The SnowSQL statements executed during the first terraform apply."
  value       = snowsql_exec.grant_all.create.0.statements
}

output "snowsql_update_stmts" {
  description = "The SnowSQL statements that will be execute as in-place changes after the first terraform apply."
  value       = snowsql_exec.grant_all.create.0.statements
}


output "snowsql_delete_stmts" {
  description = "The SnowSQL statements that will be executed during terraform destroy."
  value       = snowsql_exec.grant_all.create.0.statements
}
