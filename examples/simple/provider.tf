terraform {
  required_version = ">= 0.13.0"

  required_providers {
    snowflake = {
      source  = "chanzuckerberg/snowflake"
      version = ">= 0.25.18"
    }
    snowsql = {
      source  = "aidanmelen/snowsql"
      version = ">= 0.3.0"
    }
    random = ">= 2.1"
  }
}

provider "snowflake" {}
provider "snowsql" {}
