terraform {
  required_version = ">= 0.12.9"

  required_providers {
    snowsql = {
      source  = "aidanmelen/snowsql"
      version = ">= 0.2.0"
    }
  }
}
