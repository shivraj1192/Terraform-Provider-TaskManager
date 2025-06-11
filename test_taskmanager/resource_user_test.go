package test_taskmanager

import (
	"fmt"
	"os"
	"testing"

	"terraform-provider-taskmanager/taskmanager"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTaskmanagerUser(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	token := os.Getenv("TOKEN")

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"taskmanager": func() (*schema.Provider, error) {
				return taskmanager.Provider("dev"), nil
			},
		},
		PreCheck: func() {
			if baseURL == "" || token == "" {
				t.Fatal("BASE_URL and TOKEN must be set for acceptance tests")
			}
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTaskmanagerUserConfig(),
				Check:  testAccCheckTaskmanagerUserExists("taskmanager_user.user_new"),
			},
		},
	})
}

// Written tests for user
func testAccCheckTaskmanagerUserConfig() string {
	return fmt.Sprintf(`
resource "taskmanager_user" "user_new" {
  uname    = "AT - TASKMANAGER UNAME"
  name     = "AT - TASKMANAGER NAME"
  email    = "AT.TASKMANAGER@gmail.com"
  password = "AT - TASKMANAGER PASSWORD"
  role     = "Member"
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
