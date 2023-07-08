data "toml_file" "example" {
  input = file("${path.module}/example.toml")
}

output "toml_file_content" {
  value = jsondecode(data.toml_file.example.content_json)
}
