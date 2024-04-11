## 0.2.0 (April 11, 2024)

NOTES:

* data-source/toml_file: The `content_json` attribute is now deprecated in favour of the new `content` attribute.

FEATURES:

* data-source/toml_file: Added `content` attribute, which is a dynamically-typed attribute with the decoded TOML content as a rich type.
* function/decode: New function, available with Terraform 1.8+, to decode TOML file content. This can be used in place of the `toml_file` data source, e.g. `provider::toml::decode(file("example.toml"))`

## 0.1.0 (July 11, 2023)

FEATURES:

* data-source/toml_file: Added for parsing TOML files.
