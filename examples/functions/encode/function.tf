resource "local_file" "my_toml_file" {
  filename = "example.toml"
  content = provider::toml::encode(
    {
      version = 2
      name    = "go-toml"
      tags    = ["go", "toml"]
      section = {
        subsection = {
          items = [
            {
              include = "something"
            },
          ]
        }
      }
    }
  )
}
