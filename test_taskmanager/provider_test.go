package testtaskmanager_test

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"terraform-provider-taskmanager/taskmanager"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mitchellh/go-homedir"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testVersion = "1.0.0"

func init() {
	testAccProvider = taskmanager.Provider(testVersion)

	raw := map[string]interface{}{
		"base_url": os.Getenv("BASE_URL"),
		"token":    os.Getenv("TOKEN"),
	}

	if raw["base_url"] == "" || raw["token"] == "" {
		configPath, _ := homedir.Expand("~/.taskmanager/tf.config")
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {
			readConfigFile(configPath, raw)
		} else {
			localConfigPath := "./taskmanager/tf.config"
			if _, err := os.Stat(localConfigPath); !os.IsNotExist(err) {
				readConfigFile(localConfigPath, raw)
			}
		}
	}

	testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))

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

func readConfigFile(filePath string, raw map[string]interface{}) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")

		switch key {
		case "BASE_URL":
			raw["base_url"] = value
		case "TOKEN":
			raw["token"] = value
		}
	}
}
