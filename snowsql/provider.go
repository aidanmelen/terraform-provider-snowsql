package snowsql

import (
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is a provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Description: "The name of the Snowflake account. Can also come from the `SNOWFLAKE_ACCOUNT` environment variable.",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ACCOUNT", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for username+password authentication. Can come from the `SNOWFLAKE_USER` environment variable.",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "Password for username+password auth. Cannot be used with `browser_auth` or `private_key_path`. Can be sourced from `SNOWFLAKE_PASSWORD` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"oauth_access_token": {
				Type:          schema.TypeString,
				Description:   "Token for use with OAuth. Generating the token is left to other tools. Cannot be used with `browser_auth`, `private_key_path`, `oauth_refresh_token` or `password`. Can be sourced from `SNOWFLAKE_OAUTH_ACCESS_TOKEN` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ACCESS_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_refresh_token"},
			},
			"oauth_refresh_token": {
				Type:          schema.TypeString,
				Description:   "Token for use with OAuth. Setup and generation of the token is left to other tools. Should be used in conjunction with `oauth_client_id`, `oauth_client_secret`, `oauth_endpoint`, `oauth_redirect_url`. Cannot be used with `browser_auth`, `private_key_path`, `oauth_access_token` or `password`. Can be sourced from `SNOWFLAKE_OAUTH_REFRESH_TOKEN` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REFRESH_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_client_id": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_CLIENT_ID` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_ID", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_refresh_token", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_client_secret": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_CLIENT_SECRET` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_SECRET", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_refresh_token", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_endpoint": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_ENDPOINT` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ENDPOINT", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_refresh_token", "oauth_redirect_url"},
			},
			"oauth_redirect_url": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_REDIRECT_URL` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REDIRECT_URL", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_refresh_token"},
			},
			"browser_auth": {
				Type:          schema.TypeBool,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_USE_BROWSER_AUTH` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Description:   "Path to a private key for using keypair authentication. Cannot be used with `browser_auth`, `oauth_access_token` or `password`. Can be sourced from `SNOWFLAKE_PRIVATE_KEY_PATH` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key"},
			},
			"private_key": {
				Type:          schema.TypeString,
				Description:   "Private Key for username+private-key auth. Cannot be used with `browser_auth` or `password`. Can be sourced from `SNOWFLAKE_PRIVATE_KEY` environment variable.",
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
				Description: "Snowflake role to use for operations. If left unset, default role for user will be used. Can be sourced from the `SNOWFLAKE_ROLE` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ROLE", nil),
			},
			"region": {
				Type:        schema.TypeString,
				Description: "[Snowflake region](https://docs.snowflake.com/en/user-guide/intro-regions.html) to use.  Required if using the [legacy format for the `account` identifier](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html#format-2-legacy-account-locator-in-a-region) in the form of `<cloud_region_id>.<cloud>`. Can be sourced from the `SNOWFLAKE_REGION` environment variable.",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REGION", "us-west-2"),
			},
			"host": {
				Type:        schema.TypeString,
				Description: "Supports passing in a custom host value to the snowflake go driver for use with privatelink.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_HOST", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "Support custom port values to snowflake go driver for use with privatelink. Can be sourced from `SNOWFLAKE_PORT` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PORT", 443),
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "Support custom protocols to snowflake go driver. Can be sourced from `SNOWFLAKE_PROTOCOL` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PROTOCOL", "https"),
			},
			"warehouse": {
				Type:        schema.TypeString,
				Description: "Sets the default warehouse. Optional. Can be sourced from SNOWFLAKE_WAREHOUSE environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_WAREHOUSE", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"snowsql_exec": resourceExec(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"snowsql_query": dataSourceQuery(),
		},
		ConfigureFunc: ConfigureProvider,
	}
}

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	account := s.Get("account").(string)
	user := s.Get("username").(string)
	password := s.Get("password").(string)
	browserAuth := s.Get("browser_auth").(bool)
	privateKeyPath := s.Get("private_key_path").(string)
	privateKey := s.Get("private_key").(string)
	privateKeyPassphrase := s.Get("private_key_passphrase").(string)
	oauthAccessToken := s.Get("oauth_access_token").(string)
	region := s.Get("region").(string)
	role := s.Get("role").(string)
	oauthRefreshToken := s.Get("oauth_refresh_token").(string)
	oauthClientID := s.Get("oauth_client_id").(string)
	oauthClientSecret := s.Get("oauth_client_secret").(string)
	oauthEndpoint := s.Get("oauth_endpoint").(string)
	oauthRedirectURL := s.Get("oauth_redirect_url").(string)
	host := s.Get("host").(string)
	protocol := s.Get("protocol").(string)
	port := s.Get("port").(int)
	warehouse := s.Get("warehouse").(string)

	if oauthRefreshToken != "" {
		accessToken, err := provider.GetOauthAccessToken(oauthEndpoint, oauthClientID, oauthClientSecret, provider.GetOauthData(oauthRefreshToken, oauthRedirectURL))
		if err != nil {
			return nil, fmt.Errorf("could not retrieve access token from refresh token")
		}
		oauthAccessToken = accessToken
	}

	dsn, err := provider.DSN(
		account,
		user,
		password,
		browserAuth,
		privateKeyPath,
		privateKey,
		privateKeyPassphrase,
		oauthAccessToken,
		region,
		role,
		host,
		protocol,
		port,
		warehouse,
	)
	if err != nil {
		return nil, fmt.Errorf("could not build dsn for snowflake connection err = %w", err)
	}

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		return nil, fmt.Errorf("Could not open snowflake database err = %w", err)
	}

	return db, nil
}
