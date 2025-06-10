package test

import (
	"testing"

	"terraform-provider-taskmanager/taskmanager"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories is a map of the provider factories used for testing
var providerFactories = map[string]func() (*schema.Provider, error){
	"taskmanager": func() (*schema.Provider, error) {
		return taskmanager.Provider("test"), nil
	},
}

// TestProvider tests that the provider can be instantiated
func TestProvider(t *testing.T) {
	if err := taskmanager.Provider("test").InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// TestProviderImpl tests that the provider implementation satisfies the interface
func TestProviderImpl(t *testing.T) {
	var _ *schema.Provider = taskmanager.Provider("test")
}
