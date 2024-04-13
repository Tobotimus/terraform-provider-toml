package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pelletier/go-toml/v2"
	"strings"
)

var (
	_ function.Function = EncodeFunction{}
)

func NewEncodeFunction() function.Function {
	return EncodeFunction{}
}

type EncodeFunction struct{}

func (r EncodeFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "encode"
}

func (r EncodeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Encode a value to TOML syntax",
		MarkdownDescription: strings.Join(
			[]string{
				"Encodes a given value to string using TOML syntax.",
				"",
				"The function maps [Terraform language values](https://developer.hashicorp.com/terraform/language/expressions/types)",
				"to TOML values in the following way:",
				"",
				"| Terraform type | TOML type                                    |",
				"|----------------|----------------------------------------------|",
				"| `string`       | `String`                                     |",
				"| `number`       | `Integer` if whole number, `Float` otherwise |",
				"| `bool`         | `Boolean`                                    |",
				"| `list(...)`    | `Array`                                      |",
				"| `set(...)`     | `Array`                                      |",
				"| `tuple(...)`   | `Array`                                      |",
				"| `map(...)`     | `Table`                                      |",
				"| `object(...)`  | `Table`                                      |",
				"| Null value     | Absent from result                           |",
				"",
				"Since the TOML format cannot fully represent all Terraform language types ",
				"(and vice versa), passing the `encode` result to `decode` will not always ",
				"produce an identical value.",
			},
			"\n",
		),
		Parameters: []function.Parameter{
			function.DynamicParameter{
				Name:                "input",
				MarkdownDescription: "Terraform value to encode",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r EncodeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var dynamicArg types.Dynamic

	resp.Error = req.Arguments.Get(ctx, &dynamicArg)

	if resp.Error != nil {
		return
	}

	objToEncode := convertFromTerraformType(dynamicArg)

	encodedContent, err := toml.Marshal(objToEncode)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(
			0,
			fmt.Sprintf("The value cannot be encoded to TOML.\n\nOriginal Error: %s", err),
		)
		return
	}

	resp.Error = resp.Result.Set(ctx, types.StringValue(string(encodedContent)))
}

func convertFromTerraformType(dynamicValue attr.Value) any {
	if dynamicValue.IsNull() {
		return nil
	}
	switch value := dynamicValue.(type) {
	case types.String:
		return value.ValueString()
	case types.Int64:
		return value.ValueInt64()
	case types.Float64:
		return value.ValueFloat64()
	case types.Number:
		bigFloat := value.ValueBigFloat()
		if bigFloat.IsInt() {
			intValue, _ := bigFloat.Int64()
			return intValue
		}
		floatValue, _ := bigFloat.Float64()
		return floatValue
	case types.Bool:
		return value.ValueBool()
	case types.List:
		return convertSliceFromTerraformType(value.Elements())
	case types.Tuple:
		return convertSliceFromTerraformType(value.Elements())
	case types.Set:
		return convertSliceFromTerraformType(value.Elements())
	case types.Map:
		return convertMapFromTerraformType(value.Elements())
	case types.Object:
		return convertMapFromTerraformType(value.Attributes())
	case types.Dynamic:
		return convertFromTerraformType(value.UnderlyingValue())
	default:
		panic(
			fmt.Sprintf("Unable to convert value %v (type %T) from Terraform type", value, value),
		)
	}
}

func convertMapFromTerraformType(elements map[string]attr.Value) map[string]any {
	result := make(map[string]any, len(elements))
	for key, value := range elements {
		convertedValue := convertFromTerraformType(value)
		result[key] = convertedValue
	}
	return result
}

func convertSliceFromTerraformType(elements []attr.Value) []any {
	result := make([]any, len(elements))
	for i, value := range elements {
		convertedValue := convertFromTerraformType(value)
		result[i] = convertedValue
	}
	return result
}
