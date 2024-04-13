package provider

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const testDecodeConfig = `
output "test" {
	value = provider::toml::decode(<<EOF
[section]
string_value = "value1"
odt_value = 1979-05-27T07:32:00Z
ldt_value = 1979-05-27T07:32:00
ld_value = 1979-05-27
lt_value = 07:32:00

[[section.subsection]]
int_value = 1
float_value = 2.1

[[section.subsection]]
int_value = 2
float_value = 3.2

EOF
	)
}
`

func TestDecodeFunction(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDecodeConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"section": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"string_value": knownvalue.StringExact("value1"),
								"odt_value":    knownvalue.StringExact("1979-05-27T07:32:00Z"),
								"ldt_value":    knownvalue.StringExact("1979-05-27T07:32:00"),
								"ld_value":     knownvalue.StringExact("1979-05-27"),
								"lt_value":     knownvalue.StringExact("07:32:00"),
								"subsection": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"int_value":   knownvalue.Int64Exact(1),
										"float_value": knownvalue.Float64Exact(2.1),
									}),
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"int_value":   knownvalue.Int64Exact(2),
										"float_value": knownvalue.Float64Exact(3.2),
									}),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}
