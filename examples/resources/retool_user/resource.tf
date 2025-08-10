resource "retool_user" "example" {
  email      = "test@example.com"
  first_name = "Test"
  last_name  = "User"
  active     = true
  metadata = {
    role = "test_role"
  }
}
