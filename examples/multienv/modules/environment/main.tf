resource "local_file" "foo" {
  count    = var.enable_foo ? 1 : 0
  content  = "foo!"
  filename = "${path.root}/foo.txt"
}

resource "local_file" "bar" {
  count    = var.enable_bar ? 1 : 0
  content  = "bar!"
  filename = "${path.root}/bar.txt"
}
