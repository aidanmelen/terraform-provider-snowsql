[![release](https://github.com/aidanmelen/terraform-provider-snowsql/actions/workflows/release.yml/badge.svg)](https://github.com/aidanmelen/terraform-provider-snowsql/actions/workflows/release.yml)
[![Tests](https://github.com/aidanmelen/terraform-provider-snowsql/actions/workflows/test.yml/badge.svg)](https://github.com/aidanmelen/terraform-provider-snowsql/actions/workflows/test.yml)

# Terraform Provider SnowSQL

The `snowsql` provider allows for the management of the `create`, `read`, `update`, and `delete` lifecycles for [Snowflake](https://www.snowflake.com) objects using [SnowSQL](https://docs.snowflake.com/en/user-guide/snowsql.html).

>**Note**:  This provider is not a drop in replacement for the robust resources implemented by [terraform-provider-snowflake](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs). For example, use the `snowflake_warehouse` resource if you need to create a virtual warehouse, Use this provider when you require fine grain control of [DCL](https://www.geeksforgeeks.org/sql-ddl-dql-dml-dcl-tcl-commands/) commands or to implement Snowflake objects that are unsupported by the Snowflake provider resources.

Similiar to the [terraform-provider-shell](https://registry.terraform.io/providers/scottwinkler/shell/latest/docs); this provider

> this is a backdoor into the Terraform runtime. You can do some pretty dangerous things with this and it is up to you to make sure you don't get in trouble.
> Since this provider is rather different than most other provider, it is recommended that you at least have some familiarity with the internals of Terraform before attempting to use this provider.

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-snowsql
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory.

```shell
$ cd examples/simple
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

## Credits

see [CREDITS](CREDITS.md) for more information.
