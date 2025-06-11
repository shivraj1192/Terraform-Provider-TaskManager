package testtaskmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTaskmanagerUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTaskmanagerUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaskmanagerUserExists("taskmanager_user.user_new"),
					testAccCheckTaskmanagerUserExists("taskmanager_user.user_new1"),
				),
			},
		},
	})
}

// Written tests
func testAccCheckTaskmanagerUserConfig() string {
	return `
resource "taskmanager_user" "user_new" {
  uname    = "AT - TASKMANAGER UNAME"
  name     = "AT - TASKMANAGER NAME"
  email    = "AT.TASKMANAGER@gmail.com"
  password = "AT - TASKMANAGER PASSWORD"
  role     = "Member"
}
`
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
