resource "snowsql_exec" "function" {
  name = local.name

  create {
    statements = <<-EOT
      CREATE OR REPLACE FUNCTION ${snowflake_database.database.name}.PUBLIC.JS_FACTORIAL(f FLOAT)
        RETURNS FLOAT
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
      SHOW USER FUNCTIONS LIKE 'JS_FACTORIAL' 
        IN DATABASE ${snowflake_database.database.name};
    EOT
  }

  delete {
    statements = <<-EOT
      DROP FUNCTION IF EXISTS 
        ${snowflake_database.database.name}.PUBLIC.JS_FACTORIAL(FLOAT);
    EOT
  }
}