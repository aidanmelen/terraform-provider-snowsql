---
page_title: "Authenticating"
subcategory: ""
description: |-
  The Snowflake provider supports multiple ways to authenticate.
---

# Authentication

The Snowflake provider supports multiple ways to authenticate:

* Password
* OAuth Access Token
* Browser Auth
* Private Key

In all cases account, username, and region are required.

### Keypair Authentication Environment Variables

You should generate the public and private keys and set up environment variables.

```shell

cd ~/.ssh
openssl genrsa -out snowflake_key 4096
openssl rsa -in snowflake_key -pubout -out snowflake_key.pub
```

To export the variables into your provider:

```shell
export SNOWFLAKE_USER="..."
export SNOWFLAKE_PRIVATE_KEY_PATH="~/.ssh/snowflake_key"
```

### OAuth Access Token

If you have an OAuth access token, export these credentials as environment variables:

```shell
export SNOWFLAKE_USER='...'
export SNOWFLAKE_OAUTH_ACCESS_TOKEN='...'
```

Note that once this access token expires, you'll need to request a new one through an external application.

### Username and Password Environment Variables

If you choose to use Username and Password Authentication, export these credentials:

```shell
export SNOWFLAKE_USER='...'
export SNOWFLAKE_PASSWORD='...'
```
