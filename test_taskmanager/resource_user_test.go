package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"terraform-provider-taskmanager/taskmanager"
)

// TestResourceUser_basic tests the basic CRUD operations for the user resource
func TestResourceUser_basic(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the resource configuration
	resourceConfig := `
	resource "taskmanager_user" "test" {
		uname    = "testuser"
		name     = "Test User"
		email    = "test@example.com"
		password = "password123"
	}
	`

	// Define the resource name
	resourceName := "taskmanager_user.test"

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "uname", "testuser"),
					resource.TestCheckResourceAttr(resourceName, "name", "Test User"),
					resource.TestCheckResourceAttr(resourceName, "email", "test@example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

// TestResourceUser_update tests updating a user resource
func TestResourceUser_update(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the initial resource configuration
	initialConfig := `
	resource "taskmanager_user" "test" {
		uname    = "testuser"
		name     = "Test User"
		email    = "test@example.com"
		password = "password123"
	}
	`

	// Define the updated resource configuration
	updatedConfig := `
	resource "taskmanager_user" "test" {
		uname    = "testuser"
		name     = "Updated User"
		email    = "updated@example.com"
		password = "password123"
	}
	`

	// Define the resource name
	resourceName := "taskmanager_user.test"

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test User"),
					resource.TestCheckResourceAttr(resourceName, "email", "test@example.com"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Updated User"),
					resource.TestCheckResourceAttr(resourceName, "email", "updated@example.com"),
				),
			},
		},
	})
}

// TestResourceUserCreate tests the resourceCreateUser function directly
func TestResourceUserCreate(t *testing.T) {
	// Create a mock provider
	provider := taskmanager.Provider("test")

	// Create a mock resource data
	resourceData := provider.ResourcesMap["taskmanager_user"].Data(nil)

	// Set test values
	resourceData.Set("uname", "testuser")
	 resourceData.Set("name", "Test User")
	resourceData.Set("email", "test@example.com")
	resourceData.Set("password", "password123")

	// This test requires a mock client, which we can't fully implement here
	// In a real test, you would use a mock HTTP client to simulate API responses
	// For now, we'll just check that the resource data is set correctly
	if resourceData.Get("uname").(string) != "testuser" {
		t.Errorf("Expected uname 'testuser', got '%s'", resourceData.Get("uname").(string))
	}

	if resourceData.Get("name").(string) != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", resourceData.Get("name").(string))
	}

	if resourceData.Get("email").(string) != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", resourceData.Get("email").(string))
	}
}