output "read_results" {
  description = "The SnowSQL query result from the read statements."
  value       = try(jsondecode(nonsensitive(snowsql_exec.role_grant_all.read_results)), null)
}

output "snowsql_create_statements" {
  description = "The SnowSQL statements executed during the first terraform apply."
  value       = snowsql_exec.role_grant_all.create.0.statements
}

output "snowsql_read_statements" {
  description = "The optional SnowSQL query statements that will be execute during every terraform apply."
  value       = try(snowsql_exec.role.read.0.statements, null)
}

output "snowsql_update_statements" {
  description = "The optional SnowSQL statements that will be execute as in-place changes after the first terraform apply."
  value       = try(snowsql_exec.role_grant_all.update.0.statements, null)
}


output "snowsql_delete_statements" {
  description = "The SnowSQL statements that will be executed during terraform destroy."
  value       = snowsql_exec.role_grant_all.delete.0.statements
}
