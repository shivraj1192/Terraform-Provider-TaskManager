package test_taskmanager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"terraform-provider-taskmanager/taskmanager"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mitchellh/go-homedir"
)

var (
	testAccProvider  *schema.Provider
	testAccProviders map[string]func() (*schema.Provider, error)
	testVersion      = "1.0.0"
)

func init() {
	testAccProvider = taskmanager.Provider(testVersion)

	raw := map[string]interface{}{
		"base_url": os.Getenv("BASE_URL"),
		"token":    os.Getenv("TOKEN"),
	}

	if err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw)); err != nil {
		fmt.Println("Warning: Provider configuration in init failed:", err)
	}

	testAccProviders = map[string]func() (*schema.Provider, error){
		"taskmanager": func() (*schema.Provider, error) {
			return taskmanager.Provider(testVersion), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := taskmanager.Provider(testVersion).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("BASE_URL") == "" && os.Getenv("TOKEN") == "" {
		configPath, _ := homedir.Expand("~/.taskmanager/tf.config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatal("Either BASE_URL and TOKEN env vars must be set, or ~/.taskmanager/tf.config must exist")
		}
	}
}

func NewNotFoundErrorf(format string, a ...interface{}) error {
	return fmt.Errorf("%w %s", errors.New("Could not found"), fmt.Sprintf(format, a...))
}
