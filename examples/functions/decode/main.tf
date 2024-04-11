terraform {
  required_version = ">=1.8"

  required_providers {
    toml = {
      source  = "registry.terraform.io/tobotimus/toml"
      version = ">=0.2.0"
    }
  }
}
