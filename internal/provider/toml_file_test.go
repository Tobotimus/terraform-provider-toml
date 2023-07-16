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
	testAccTomlFileDataSourceExpectedOutput = `
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
	if err := json.Compact(dst, []byte(testAccTomlFileDataSourceExpectedOutput)); err != nil {
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.toml_file.file", "content_json", expected_content_minified),
					resource.TestCheckResourceAttr("data.toml_file.file", "id", sha1Hex),
				),
			},
		},
	})
}
