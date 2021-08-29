package snowsql

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"snowsql": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("SNOWFLAKE_USERNAME"); err == "" {
		t.Fatal("SNOWFLAKE_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("SNOWFLAKE_PASSWORD"); err == "" {
		t.Fatal("SNOWFLAKE_PASSWORD must be set for acceptance tests")
	}
}
