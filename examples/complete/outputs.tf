output "snowsql_create_statements" {
  description = "The SnowSQL statements executed during the first terraform apply."
  value       = snowsql_exec.grant_all.create.0.statements
}

output "snowsql_update_statements" {
  description = "The SnowSQL statements that will be execute as in-place changes after the first terraform apply."
  value       = try(snowsql_exec.grant_all.update.0.statements, null)
}


output "snowsql_delete_statements" {
  description = "The SnowSQL statements that will be executed during terraform destroy."
  value       = snowsql_exec.grant_all.delete.0.statements
}
