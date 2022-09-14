resource "google_project_iam_policy" "project" {
  project     = "jond-default"
  policy_data = data.google_iam_policy.policy.policy_data
}

data "google_iam_policy" "policy" {
  binding {
    role = "roles/owner"
    members = [
      "user:jond3k+missing@gmail.com",
    ]
  }
}
