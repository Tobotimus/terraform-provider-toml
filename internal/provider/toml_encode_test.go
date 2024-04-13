package provider

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	testEncodeConfig = `
output "test" {
	value = provider::toml::encode({
		"section": {
			"int": 1,
			"float": 2.1,
			"subsection": [{"string": "value"}],
		},
		"another_section": {
			"boolean": true,
			"null_value": null,
			"set_value": toset(["b", "a"]),
		},
	})
}
`

	testEncodeExpectedOutput = `[another_section]
boolean = true
set_value = ['a', 'b']

[section]
float = 2.1
int = 1

[[section.subsection]]
string = 'value'
`
)

func TestEncodeFunction(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testEncodeConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact(testEncodeExpectedOutput),
					),
				},
			},
		},
	})
}
