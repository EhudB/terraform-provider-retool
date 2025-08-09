package user_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/tryretool/terraform-provider-retool/internal/acctest"
)

const testActiveUser = `
	resource "retool_user" "test_user" {
		email = "test@example.com"
		first_name = "Test"
		last_name = "User"
		active = true
		metadata = {
			"role": "test_role"
		}
	}
		`

const testUpdatedUser = `
	resource "retool_user" "test_user_updated" {
		email = "updated_test@example.com"
		first_name = "Updated"
		last_name = "User"
		active = true
		metadata = {
			"role": "updated_role"
		}
	}
		`


const testDisabeldUser = `
	resource "retool_user" "test_user_disabled" {
		email = "test@example.com"
		first_name = "Test"
		last_name = "User"
		active = false
		metadata = {
			"role": "updated_role"
		}
	}
		`

const testDefaultValuesConfig = `
	resource "retool_user" "test_user_with_defaults" {
		email = "test@example.com"
		first_name = "Test"
		last_name = "User"
	}
		`

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func TestAccUser(t *testing.T) {
	acctest.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			// Read and Create.
			{
				Config: testActiveUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("retool_user.test_user", "email", "test@example.com"),
					resource.TestCheckResourceAttr("retool_user.test_user", "first_name", "Test"),
					resource.TestCheckResourceAttr("retool_user.test_user", "last_name", "User"),
					resource.TestCheckResourceAttr("retool_user.test_user", "active", "true"),
					resource.TestCheckResourceAttr("retool_user.test_user", "metadata.role", "test_role"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "id"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "legacy_id"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "created_at"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "last_active"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "user_type"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "two_factor_auth_enabled"),
				),
			},
			// Import state.
			{
				ResourceName:      "retool_user.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read.
			{
				Config: testUpdatedUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("retool_user.test_user", "email", "updated_test@example.com"),
					resource.TestCheckResourceAttr("retool_user.test_user", "first_name", "Updated"),
					resource.TestCheckResourceAttr("retool_user.test_user", "last_name", "User"),
					resource.TestCheckResourceAttr("retool_user.test_user", "active", "true"),
					resource.TestCheckResourceAttr("retool_user.test_user", "metadata.role", "updated_role"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "id"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "legacy_id"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "created_at"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "last_active"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "user_type"),
					resource.TestCheckResourceAttrSet("retool_user.test_user", "two_factor_auth_enabled"),
				),
			},
			// Check default values.
			{
				Config: testDefaultValuesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("retool_user.test_user_with_defaults", "email", "test@example.com"),
					resource.TestCheckResourceAttr("retool_user.test_user_with_defaults", "first_name", "Test"),
					resource.TestCheckResourceAttr("retool_user.test_user_with_defaults", "last_name", "User"),
					resource.TestCheckResourceAttr("retool_user.test_user_with_defaults", "active", "true"),
					resource.TestCheckResourceAttrSet("retool_user.test_user_with_defaults", "id"),
					resource.TestCheckResourceAttrSet("retool_user.test_user_with_defaults", "legacy_id"),
					resource.TestCheckResourceAttrSet("retool_user.test_user_with_defaults", "created_at"),
					resource.TestCheckResourceAttrSet("retool_user.test_user_with_defaults", "last_active"),
					resource.TestCheckResourceAttrSet("retool_user.test_user_with_defaults", "user_type"),
					resource.TestCheckResourceAttrSet("retool_user.test_user_with_defaults", "two_factor_auth_enabled"),
				),
			},
		},
	})
}

func sweepUsers(region string) error {
	log.Printf("Sweeping users in region %s", region)
	client, err := acctest.SweeperClient()
	if err != nil {
		return err
	}

	users, _, err := client.UsersAPI.UsersGet(context.Background()).Execute()
	if err != nil {
		return fmt.Errorf("Error reading users: %s", err.Error())
	}

	for _, user := range users.Data {
		if strings.HasPrefix(user.Email, "tf-acc") {
			log.Printf("Deleting user %s", user.Email)
			_, err := client.UsersAPI.UsersUserIdDelete(context.Background(), user.Id).Execute()
			if err != nil {
				return fmt.Errorf("Error deleting user %s: %s", user.Email, err.Error())
			}
		}
	}
	return nil
}

func init() {
	resource.AddTestSweepers("retool_user", &resource.Sweeper{
		Name: "retool_user",
		F:    sweepUsers,
	})
}
