package test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"terraform-provider-taskmanager/taskmanager"
)

// testAccPreCheck validates the necessary test API credentials exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TASKMANAGER_TOKEN"); v == "" {
		t.Fatal("TASKMANAGER_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("TASKMANAGER_BASE_URL"); v == "" {
		t.Fatal("TASKMANAGER_BASE_URL must be set for acceptance tests")
	}
}

// testAccProvider is the "main" provider instance
func testAccProvider() *schema.Provider {
	return taskmanager.Provider("test")
}

// testAccProviderFactories is a static map containing only the main provider instance
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"taskmanager": func() (*schema.Provider, error) {
		return testAccProvider(), nil
	},
}

// testAccProviderConfigure configures the provider with test credentials
func testAccProviderConfigure(t *testing.T) {
	provider := testAccProvider()

	raw := map[string]interface{}{
		"token":    os.Getenv("TASKMANAGER_TOKEN"),
		"base_url": os.Getenv("TASKMANAGER_BASE_URL"),
	}

	err := provider.Configure(nil, terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatal(err)
	}
}

// skipIfShortTest skips a test if running in short mode
func skipIfShortTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}
}