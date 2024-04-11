// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const (
	testAccTomlFileDataSourceConfig = `
		data "toml_file" "file" {
		  input = <<EOF
		version = 2
		name = "go-toml"
		tags = ["go", "toml"]

		[section.subsection]
		items = [
			{include = "something"}
		]
		EOF

		}
	`

	// Note: Keys MUST be sorted alphabetically.
	testAccTomlFileDataSourceExpectedOutputJSON = `
		{
          "name": "go-toml",
          "section": {
            "subsection": {
              "items": [
                {
                  "include": "something"
                }
              ]
            }
          },
          "tags": [
            "go",
            "toml"
          ],
		  "version": 2
		}
	`
)

func TestAccTomlFileDataSource(t *testing.T) {
	dst := &bytes.Buffer{}
	if err := json.Compact(dst, []byte(testAccTomlFileDataSourceExpectedOutputJSON)); err != nil {
		panic(err)
	}
	expected_content_minified := dst.String()

	sha1Sum := sha1.Sum([]byte(expected_content_minified))
	sha1Hex := hex.EncodeToString(sha1Sum[:])

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: testAccTomlFileDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.toml_file.file",
						tfjsonpath.New("content"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("go-toml"),
							"section": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"subsection": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"items": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"include": knownvalue.StringExact("something"),
										}),
									}),
								}),
							}),
							"tags": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("go"),
								knownvalue.StringExact("toml"),
							}),
							"version": knownvalue.Int64Exact(2),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.toml_file.file",
						tfjsonpath.New("content_json"),
						knownvalue.StringExact(expected_content_minified),
					),
					statecheck.ExpectKnownValue(
						"data.toml_file.file",
						tfjsonpath.New("id"),
						knownvalue.StringExact(sha1Hex),
					),
				},
			},
		},
	})
}
