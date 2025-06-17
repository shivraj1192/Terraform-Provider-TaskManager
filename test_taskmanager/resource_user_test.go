package test_taskmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTaskmanager(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTaskmanagerUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaskmanagerUserExists("taskmanager_user.user_new"),
					testAccCheckTaskmanagerUserExists("taskmanager_team.team_new"),
					testAccCheckTaskmanagerUserExists("taskmanager_task.task_new"),
					testAccCheckTaskmanagerUserExists("taskmanager_comment.comment_new"),
					testAccCheckTaskmanagerUserExists("taskmanager_attachment.attachment_new"),
				),
			},
		},
	})
}

func testAccCheckTaskmanagerUserConfig() string {
	return fmt.Sprintf(`
resource "taskmanager_user" "user_new" {
  uname    = "AT - TASKMANAGER UNAME"
  name     = "AT - TASKMANAGER NAME"
  email    = "AT.TASKMANAGER@gmail.com"
  password = "AT - TASKMANAGER PASSWORD"
  role     = "Member"
}

resource "taskmanager_team" "team_new" {
  name        = "AT - TASKMANAGER NAME"
  description = "AT - TASKMANAGER DESC"
  members     = [1, taskmanager_user.user_new.id]
}

resource "taskmanager_task" "task_new" {
  title = "AT - TASKMANAGER TITLE"
  description = "AT - TASKMANAGER DESC"
  priority = "High"
  status = "In progress"
  due_date = "2025-06-25T18:30:00Z"
  team_id = taskmanager_team.team_new.id
  parent_task_id = 0
  assignees = [taskmanager_user.user_new.id]
  labels = []
}

resource "taskmanager_comment" "comment_new"{
  content = "AT - TASKMANAGER CONTENT"
  task_id = taskmanager_task.task_new.id
  parent_comment_id = 0
}

resource "taskmanager_attachment" "attachment_new"{
  file_name = "infralovers_courses.pdf"
  task_id = taskmanager_task.task_new.id
  url = "./static/files/infralovers_courses.pdf"
}
`)
}

func testAccCheckTaskmanagerUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID set for resource: %s", resourceName)
		}
		return nil
	}
}
