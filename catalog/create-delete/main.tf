resource "local_file" "foo" {
  content  = "bar!"
  filename = "${path.module}/foo.bar"
  lifecycle {
    create_before_destroy = true
  }
}
