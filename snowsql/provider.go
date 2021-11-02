package snowsql

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
	"github.com/youmark/pkcs8"
	"golang.org/x/crypto/ssh"
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
				ConflictsWith: []string{"browser_auth", "private_key_path", "oauth_access_token"},
			},
			"oauth_access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ACCESS_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "password"},
			},
			"browser_auth": {
				Type:          schema.TypeBool,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password", "private_key_path", "oauth_access_token"},
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token"},
			},
			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key_path"},
			},
			"private_key_passphrase": {
				Type:          schema.TypeString,
				Description:   "Supports the encryption ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PASSPHRASE", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token"},
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
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

//providerConfigure returns an authenticated Snowflake connection.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	account := d.Get("account").(string)
	user := d.Get("username").(string)
	password := d.Get("password").(string)
	browserAuth := d.Get("browser_auth").(bool)
	privateKeyPath := d.Get("private_key_path").(string)
	privateKey := d.Get("private_key").(string)
	privateKeyPassphrase := d.Get("private_key_passphrase").(string)
	oauthAccessToken := d.Get("oauth_access_token").(string)
	region := d.Get("region").(string)
	role := d.Get("role").(string)

	var diags diag.Diagnostics

	dsn, err := DSN(account, user, password, browserAuth, privateKeyPath, privateKey, privateKeyPassphrase, oauthAccessToken, region, role)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not build data source name (DSN) for Snowflake connection",
		})
		return nil, diags
	}

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not open Snowflake database.",
		})

		return nil, diags
	}

	return db, diags
}

// DSN returns a Snowflake Data Source Name used to authenticate the provider.
func DSN(
	account,
	user,
	password string,
	browserAuth bool,
	privateKeyPath,
	privateKey,
	privateKeyPassphrase,
	oauthAccessToken,
	region,
	role string) (string, error) {

	// us-west-2 is their default region, but if you actually specify that it won't trigger their default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if region == "us-west-2" {
		region = ""
	}

	// https://godoc.org/github.com/snowflakedb/gosnowflake#Config
	config := gosnowflake.Config{
		Account: account,
		User:    user,
		Region:  region,
		Role:    role,
	}

	if privateKeyPath != "" {
		privateKeyBytes, err := ReadPrivateKeyFile(privateKeyPath)
		if err != nil {
			return "", errors.Wrap(err, "Private Key file could not be read")
		}
		rsaPrivateKey, err := ParsePrivateKey(privateKeyBytes, []byte(privateKeyPassphrase))
		if err != nil {
			return "", errors.Wrap(err, "Private Key could not be parsed")
		}
		config.PrivateKey = rsaPrivateKey
		config.Authenticator = gosnowflake.AuthTypeJwt
	} else if privateKey != "" {
		rsaPrivateKey, err := ParsePrivateKey([]byte(privateKey), []byte(privateKeyPassphrase))
		if err != nil {
			return "", errors.Wrap(err, "Private Key could not be parsed")
		}
		config.PrivateKey = rsaPrivateKey
		config.Authenticator = gosnowflake.AuthTypeJwt
	} else if browserAuth {
		config.Authenticator = gosnowflake.AuthTypeExternalBrowser
	} else if oauthAccessToken != "" {
		config.Authenticator = gosnowflake.AuthTypeOAuth
		config.Token = oauthAccessToken
	} else if password != "" {
		config.Password = password
	} else {
		return "", errors.New("no authentication method provided")
	}

	return gosnowflake.DSN(&config)
}

// ParsePrivateKey reads and parses an RSA Private Key.
func ParsePrivateKey(privateKeyBytes []byte, passphrase []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("Could not parse private key, key is not in PEM format")
	}

	if privateKeyBlock.Type == "ENCRYPTED PRIVATE KEY" {
		if len(passphrase) == 0 {
			return nil, fmt.Errorf("Private key requires a passphrase, but private_key_passphrase was not supplied")
		}
		privateKey, err := pkcs8.ParsePKCS8PrivateKeyRSA(privateKeyBlock.Bytes, passphrase)
		if err != nil {
			return nil, errors.Wrap(
				err,
				"Could not parse encrypted private key with passphrase, only ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc are supported",
			)
		}
		return privateKey, nil
	}

	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse private key")
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("privateKey not of type RSA")
	}
	return rsaPrivateKey, nil
}


func ReadPrivateKeyFile(privateKeyPath string) ([]byte, error) {
	expandedPrivateKeyPath, err := homedir.Expand(privateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Path to private key")
	}

	privateKeyBytes, err := ioutil.ReadFile(expandedPrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read private key")
	}

	if len(privateKeyBytes) == 0 {
		return nil, errors.New("Private key is empty")
	}

	return privateKeyBytes, nil
}