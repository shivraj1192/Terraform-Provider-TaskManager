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
				),
			},
		},
	})
}

func testAccCheckTaskmanagerUserConfig() string {
	return fmt.Sprintf(`
	resource "taskmanager_user" "user_new" {
	uname = "AT - TASKMANAGER UNAME"
	name = "AT - TASKMANAGER NAME"
	email = "AT.TASKMANAGER@gmail.com"
	password = "AT - TASKMANAGER PASSWORD"
	role = "Member"
	}

	resource "taskmanager_user" "user_new1" {
	uname = "AT - TASKMANAGER UNAME1"
	name = "AT - TASKMANAGER NAME1"
	email = "AT.TASKMANAGER1@gmail.com"
	password = "AT - TASKMANAGER PASSWORD1"
	role = "Member"
	}
	`)
}

func testAccCheckTaskmanagerUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return NewNotFoundErrorf("%s in state", n)
		}

		if rs.Primary.ID == "" {
			return NewNotFoundErrorf("ID for %s in state", n)
		}

		return nil
	}
}
