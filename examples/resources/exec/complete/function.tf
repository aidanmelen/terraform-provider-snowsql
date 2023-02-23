resource "snowsql_exec" "function" {
  name = local.name

  create {
    statements = <<-EOT
      USE SCHEMA ${snowflake_database.database.name}.PUBLIC;

      CREATE OR REPLACE FUNCTION js_factorial(d double)
        RETURNS double
        LANGUAGE JAVASCRIPT
        STRICT
        AS '
        if (D <= 0) {
          return 1;
        } else {
          var result = 1;
          for (var i = 2; i <= D; i++) {
            result = result * i;
          }
          return result;
        }
        ';
    EOT
  }

  read {
    statements = <<-EOT
      SHOW USER FUNCTIONS LIKE 'js_factorial' 
        IN DATABASE ${snowflake_database.database.name};
    EOT
  }

  delete {
    statements = <<-EOT
      DROP FUNCTION IF EXISTS 
        ${snowflake_database.database.name}.PUBLIC.js_factorial(FLOAT);
    EOT
  }
}