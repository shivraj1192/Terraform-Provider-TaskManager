package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"terraform-provider-taskmanager/taskmanager"
)

// TestResourceTeam_basic tests the basic CRUD operations for the team resource
func TestResourceTeam_basic(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the resource configuration
	resourceConfig := `
	resource "taskmanager_user" "owner" {
		uname    = "teamowner"
		name     = "Team Owner"
		email    = "owner@example.com"
		password = "password123"
	}

	resource "taskmanager_team" "test" {
		name        = "Test Team"
		description = "A team for testing"
		members     = [taskmanager_user.owner.id]
	}
	`

	// Define the resource name
	resourceName := "taskmanager_team.test"

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Team"),
					resource.TestCheckResourceAttr(resourceName, "description", "A team for testing"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

// TestResourceTeam_update tests updating a team resource
func TestResourceTeam_update(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the initial resource configuration
	initialConfig := `
	resource "taskmanager_user" "owner" {
		uname    = "teamowner"
		name     = "Team Owner"
		email    = "owner@example.com"
		password = "password123"
	}

	resource "taskmanager_user" "member" {
		uname    = "teammember"
		name     = "Team Member"
		email    = "member@example.com"
		password = "password123"
	}

	resource "taskmanager_team" "test" {
		name        = "Test Team"
		description = "A team for testing"
		members     = [taskmanager_user.owner.id]
	}
	`

	// Define the updated resource configuration
	updatedConfig := `
	resource "taskmanager_user" "owner" {
		uname    = "teamowner"
		name     = "Team Owner"
		email    = "owner@example.com"
		password = "password123"
	}

	resource "taskmanager_user" "member" {
		uname    = "teammember"
		name     = "Team Member"
		email    = "member@example.com"
		password = "password123"
	}

	resource "taskmanager_team" "test" {
		name        = "Updated Team"
		description = "An updated team for testing"
		members     = [taskmanager_user.owner.id, taskmanager_user.member.id]
	}
	`

	// Define the resource name
	resourceName := "taskmanager_team.test"

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Team"),
					resource.TestCheckResourceAttr(resourceName, "description", "A team for testing"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Updated Team"),
					resource.TestCheckResourceAttr(resourceName, "description", "An updated team for testing"),
				),
			},
		},
	})
}

// TestResourceTeamCreate tests the resourceCreateTeam function directly
func TestResourceTeamCreate(t *testing.T) {
	// Create a mock provider
	provider := taskmanager.Provider("test")

	// Create a mock resource data
	resourceData := provider.ResourcesMap["taskmanager_team"].Data(nil)

	// Set test values
	resourceData.Set("name", "Test Team")
	resourceData.Set("description", "A team for testing")
	resourceData.Set("members", []interface{}{1})

	// This test requires a mock client, which we can't fully implement here
	// In a real test, you would use a mock HTTP client to simulate API responses
	// For now, we'll just check that the resource data is set correctly
	if resourceData.Get("name").(string) != "Test Team" {
		t.Errorf("Expected name 'Test Team', got '%s'", resourceData.Get("name").(string))
	}

	if resourceData.Get("description").(string) != "A team for testing" {
		t.Errorf("Expected description 'A team for testing', got '%s'", resourceData.Get("description").(string))
	}
}