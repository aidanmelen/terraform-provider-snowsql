package snowsql

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is a provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ACCOUNT", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"oauth_access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ACCESS_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_refresh_token"},
			},
			"oauth_refresh_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REFRESH_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_ID", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_refresh_token", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_SECRET", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_refresh_token", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_endpoint": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ENDPOINT", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_refresh_token", "oauth_redirect_url"},
			},
			"oauth_redirect_url": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REDIRECT_URL", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_refresh_token"},
			},
			"browser_auth": {
				Type:          schema.TypeBool,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key"},
			},
			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key_path", "oauth_refresh_token"},
			},
			"private_key_passphrase": {
				Type:          schema.TypeString,
				Description:   "Supports the encryption ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PASSPHRASE", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "oauth_refresh_token"},
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ROLE", nil),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REGION", "us-west-2"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"snowsql_exec": resourceExec(),
		},
		ConfigureFunc: provider.ConfigureProvider,
	}
}
