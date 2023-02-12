output "results" {
  description = "The SnowSQL query result from the read statements."
  value       = try(snowsql_exec.role, null)
}

output "snowsql_create_statements" {
  description = "The SnowSQL statements executed during the first terraform apply."
  value       = snowsql_exec.role.create.0.statements
}

output "snowsql_read_statements" {
  description = "The SnowSQL query statements that will be execute during every terraform apply."
  value       = try(snowsql_exec.role.read.0.statements, null)
}

output "snowsql_update_statements" {
  description = "The SnowSQL statements that will be execute as in-place changes after the first terraform apply."
  value       = try(snowsql_exec.role.update.0.statements, null)
}


output "snowsql_delete_statements" {
  description = "The SnowSQL statements that will be executed during terraform destroy."
  value       = snowsql_exec.role.delete.0.statements
}

