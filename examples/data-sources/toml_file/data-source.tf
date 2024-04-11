data "toml_file" "example" {
  input = file("${path.module}/example.toml")
}

output "toml_file_content" {
  value = data.toml_file.example.content
}
