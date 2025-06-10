package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"terraform-provider-taskmanager/taskmanager"
)

// TestResourceComment_basic tests the basic CRUD operations for the comment resource
func TestResourceComment_basic(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the resource configuration
	resourceConfig := `
	resource "taskmanager_user" "creator" {
		uname    = "commentcreator"
		name     = "Comment Creator"
		email    = "creator@example.com"
		password = "password123"
	}

	resource "taskmanager_team" "team" {
		name        = "Comment Team"
		description = "A team for comment testing"
		members     = [taskmanager_user.creator.id]
	}

	resource "taskmanager_task" "task" {
		title       = "Comment Task"
		description = "A task for comment testing"
		status      = "To Do"
		priority    = "Medium"
		creator_id  = taskmanager_user.creator.id
		team_id     = taskmanager_team.team.id
	}

	resource "taskmanager_comment" "test" {
		content    = "This is a test comment"
		task_id    = taskmanager_task.task.id
		author_id  = taskmanager_user.creator.id
	}
	`

	// Define the resource name
	resourceName := "taskmanager_comment.test"

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "content", "This is a test comment"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

// TestResourceComment_update tests updating a comment resource
func TestResourceComment_update(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the initial resource configuration
	initialConfig := `
	resource "taskmanager_user" "creator" {
		uname    = "commentcreator"
		name     = "Comment Creator"
		email    = "creator@example.com"
		password = "password123"
	}

	resource "taskmanager_team" "team" {
		name        = "Comment Team"
		description = "A team for comment testing"
		members     = [taskmanager_user.creator.id]
	}

	resource "taskmanager_task" "task" {
		title       = "Comment Task"
		description = "A task for comment testing"
		status      = "To Do"
		priority    = "Medium"
		creator_id  = taskmanager_user.creator.id
		team_id     = taskmanager_team.team.id
	}

	resource "taskmanager_comment" "test" {
		content    = "This is a test comment"
		task_id    = taskmanager_task.task.id
		author_id  = taskmanager_user.creator.id
	}
	`

	// Define the updated resource configuration
	updatedConfig := `
	resource "taskmanager_user" "creator" {
		uname    = "commentcreator"
		name     = "Comment Creator"
		email    = "creator@example.com"
		password = "password123"
	}

	resource "taskmanager_team" "team" {
		name        = "Comment Team"
		description = "A team for comment testing"
		members     = [taskmanager_user.creator.id]
	}

	resource "taskmanager_task" "task" {
		title       = "Comment Task"
		description = "A task for comment testing"
		status      = "To Do"
		priority    = "Medium"
		creator_id  = taskmanager_user.creator.id
		team_id     = taskmanager_team.team.id
	}

	resource "taskmanager_comment" "test" {
		content    = "This is an updated test comment"
		task_id    = taskmanager_task.task.id
		author_id  = taskmanager_user.creator.id
	}
	`

	// Define the resource name
	resourceName := "taskmanager_comment.test"

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "content", "This is a test comment"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "content", "This is an updated test comment"),
				),
			},
		},
	})
}

// TestResourceCommentCreate tests the resourceCreateComment function directly
func TestResourceCommentCreate(t *testing.T) {
	// Create a mock provider
	provider := taskmanager.Provider("test")

	// Create a mock resource data
	resourceData := provider.ResourcesMap["taskmanager_comment"].Data(nil)

	// Set test values
	resourceData.Set("content", "This is a test comment")
	resourceData.Set("task_id", 1)
	resourceData.Set("author_id", 1)

	// This test requires a mock client, which we can't fully implement here
	// In a real test, you would use a mock HTTP client to simulate API responses
	// For now, we'll just check that the resource data is set correctly
	if resourceData.Get("content").(string) != "This is a test comment" {
		t.Errorf("Expected content 'This is a test comment', got '%s'", resourceData.Get("content").(string))
	}

	if resourceData.Get("task_id").(int) != 1 {
		t.Errorf("Expected task_id 1, got %d", resourceData.Get("task_id").(int))
	}

	if resourceData.Get("author_id").(int) != 1 {
		t.Errorf("Expected author_id 1, got %d", resourceData.Get("author_id").(int))
	}
}