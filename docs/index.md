# SnowSQL Provider

The Terraform SnowSQL provider allows for the management of the `create` and `delete` lifecycles for [Snowflake](https://www.snowflake.com) objects with [SnowSQL](https://docs.snowflake.com/en/user-guide/snowsql.html).

## Example

```hcl
terraform {
  required_providers {
    snowsql = {
      source  = "aidanmelen/snowsql"
      version = ">= 0.1.0"
    }
  }
}

provider snowsql {
  // required
  username = "..."
  account  = "..."
  region   = "..."

  // optional, at exactly one must be set
  password           = "..."
  oauth_access_token = "..."
  private_key_path   = "..."

  // optional
  role = "..."
}
```

The SnowSQL provider is intended to be used in conjunction with the Snowflake Provider. In fact, both providers share the same authorization. This means we can reuse the authentication variables to configure both providers. For example:

```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "chanzuckerberg/snowflake"
      version = ">= v0.25.18"
    }
    snowsql = {
      source  = "aidanmelen/snowsql"
      version = ">= 0.1.0"
    }
  }
}

# assuming environment variables are configured
provider "snowflake" {}
provider "snowsql" {}
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Snowflake
 `provider` block:

* `account` - (required) The name of the Snowflake account. Can also come from the
  `SNOWFLAKE_ACCOUNT` environment variable.
* `username` - (required) Username for username+password authentication. Can come from the
  `SNOWFLAKE_USER` environment variable.
* `region` - (required) [Snowflake region](https://docs.snowflake.com/en/user-guide/intro-regions.html) to use. Can be source from the `SNOWFLAKE_REGION` environment variable.
* `password` - (optional) Password for username+password auth. Cannot be used with `browser_auth` or
  `private_key_path`. Can be source from `SNOWFLAKE_PASSWORD` environment variable.
* `oauth_access_token` - (optional) Token for use with OAuth. Generating the token is left to other
  tools. Cannot be used with `browser_auth`, `private_key_path` or `password`. Can be source from
  `SNOWFLAKE_OAUTH_ACCESS_TOKEN` environment variable.
* `private_key_path` - (optional) Path to a private key for using keypair authentication.. Cannot be
  used with `browser_auth`, `oauth_access_token` or `password`. Can be source from
  `SNOWFLAKE_PRIVATE_KEY_PATH` environment variable.
* `role` - (optional) Snowflake role to use for operations. If left unset, default role for user
  will be used. Can come from the `SNOWFLAKE_ROLE` environment variable.
