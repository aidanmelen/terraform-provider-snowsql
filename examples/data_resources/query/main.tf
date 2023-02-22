data "snowsql_query" "accountadmin_role" {
  statements = "SHOW ROLES LIKE 'ACCOUNTADMIN'"
}