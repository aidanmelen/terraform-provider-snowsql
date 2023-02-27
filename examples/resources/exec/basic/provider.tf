terraform {
  required_version = ">= 0.13.0"

  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = ">= 0.56.5"
    }
    snowsql = {
      source  = "aidanmelen/snowsql"
      version = ">= 1.3.2"
    }
    random = ">= 2.1"
  }
}

provider "snowflake" {}
provider "snowsql" {}
