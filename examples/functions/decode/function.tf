output "toml_file_content" {
  value = provider::toml::decode(file("${path.module}/example.toml"))
}
