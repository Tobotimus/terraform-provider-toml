## 0.3.1 (July 15, 2024)

NOTES:

* Minor documentation improvements. 

## 0.3.0 (April 13, 2024)

FEATURES:

* function/encode: New function, available with Terraform 1.8+, to encode a value as a string using TOML syntax.

BUG FIXES:

* function/decode: Fix errors when decoding date-times and dates.

NOTES:

* function/decode: Documentation updated to include type mapping between TOML values and Terraform language values.

## 0.2.1 (April 11, 2024)

NOTES:

* data-source/toml_file: Update documentation to remove recommendation to use `content_json` attribute.


## 0.2.0 (April 11, 2024)

NOTES:

* data-source/toml_file: The `content_json` attribute is now deprecated in favour of the new `content` attribute.

FEATURES:

* data-source/toml_file: Added `content` attribute, which is a dynamically-typed attribute with the decoded TOML content as a rich type.
* function/decode: New function, available with Terraform 1.8+, to decode TOML file content. This can be used in place of the `toml_file` data source, e.g. `provider::toml::decode(file("example.toml"))`

## 0.1.0 (July 11, 2023)

FEATURES:

* data-source/toml_file: Added for parsing TOML files.
