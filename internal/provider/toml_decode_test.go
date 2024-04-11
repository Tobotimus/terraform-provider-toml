package provider

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const testConfig = `
output "test" {
	value = provider::toml::decode(<<EOF
[section]
field1 = "value1"
EOF
	)
}
`

func TestDecodeFunction_Known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"section": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"field1": knownvalue.StringExact("value1"),
							}),
						}),
					),
				},
			},
		},
	})
}
