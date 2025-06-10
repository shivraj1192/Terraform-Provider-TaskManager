package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestDataSourceUser tests the user data source
func TestDataSourceUser(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the resource and data source configuration
	config := `
	resource "taskmanager_user" "test" {
		uname    = "datasourceuser"
		name     = "Data Source User"
		email    = "datasource@example.com"
		password = "password123"
	}

	data "taskmanager_user" "test" {
		id = taskmanager_user.test.id
	}
	`

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.taskmanager_user.test", "name",
						"taskmanager_user.test", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.taskmanager_user.test", "email",
						"taskmanager_user.test", "email",
					),
				),
			},
		},
	})
}

// TestDataSourceTeam tests the team data source
func TestDataSourceTeam(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the resource and data source configuration
	config := `
	resource "taskmanager_user" "owner" {
		uname    = "teamowner"
		name     = "Team Owner"
		email    = "owner@example.com"
		password = "password123"
	}

	resource "taskmanager_team" "test" {
		name        = "Data Source Team"
		description = "A team for data source testing"
		members     = [taskmanager_user.owner.id]
	}

	data "taskmanager_team" "test" {
		id = taskmanager_team.test.id
	}
	`

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.taskmanager_team.test", "name",
						"taskmanager_team.test", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.taskmanager_team.test", "description",
						"taskmanager_team.test", "description",
					),
				),
			},
		},
	})
}

// TestDataSourceTask tests the task data source
func TestDataSourceTask(t *testing.T) {
	// Skip this test in short mode as it requires a running API
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	// Define the resource and data source configuration
	config := `
	resource "taskmanager_user" "creator" {
		uname    = "taskcreator"
		name     = "Task Creator"
		email    = "creator@example.com"
		password = "password123"
	}

	resource "taskmanager_team" "team" {
		name        = "Task Team"
		description = "A team for task testing"
		members     = [taskmanager_user.creator.id]
	}

	resource "taskmanager_task" "test" {
		title       = "Data Source Task"
		description = "A task for data source testing"
		status      = "To Do"
		priority    = "Medium"
		due_date    = "2023-12-31T23:59:59Z"
		creator_id  = taskmanager_user.creator.id
		team_id     = taskmanager_team.team.id
		assignees   = [taskmanager_user.creator.id]
	}

	data "taskmanager_task" "test" {
		id = taskmanager_task.test.id
	}
	`

	// Run the acceptance test
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.taskmanager_task.test", "title",
						"taskmanager_task.test", "title",
					),
					resource.TestCheckResourceAttrPair(
						"data.taskmanager_task.test", "description",
						"taskmanager_task.test", "description",
					),
					resource.TestCheckResourceAttrPair(
						"data.taskmanager_task.test", "status",
						"taskmanager_task.test", "status",
					),
				),
			},
		},
	})
}
