package testtaskmanager_test

import (
	"errors"
	"fmt"
	"os"
	"terraform-provider-taskmanager/taskmanager"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testVersion = "1.0.0"

func init() {
	testAccProvider = taskmanager.Provider(testVersion)
	testAccProviders = map[string]*schema.Provider{
		"taskmanager": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := taskmanager.Provider(testVersion).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = taskmanager.Provider(testVersion)
}

func testAccPreCheck(t *testing.T) {
	configPath, _ := homedir.Expand("~/.taskmanager/tf.config")
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		return
	}
	if err := os.Getenv("BASE_URL"); err == "" {
		t.Fatal("BASE_URL must be set for acceptance tests")
	}
	if err := os.Getenv("TOKEN"); err == "" {
		t.Fatal("TOKEN must be set for acceptance tests")
	}
}

func NewNotFoundErrorf(format string, a ...interface{}) error {
	return fmt.Errorf("%w %s", errors.New("Could not found"), fmt.Sprintf(format, a...))
}
