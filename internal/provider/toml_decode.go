package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pelletier/go-toml/v2"
)

var (
	_ function.Function = DecodeFunction{}
)

func NewDecodeFunction() function.Function {
	return DecodeFunction{}
}

type DecodeFunction struct{}

func (r DecodeFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "decode"
}

func (r DecodeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Decode TOML content",
		MarkdownDescription: "Decodes the content of a TOML file to a Terraform object.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "input",
				MarkdownDescription: "TOML file content to decode",
			},
		},
		Return: function.DynamicReturn{},
	}
}

func (r DecodeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data string

	resp.Error = req.Arguments.Get(ctx, &data)

	if resp.Error != nil {
		return
	}

	var decoded_content any
	err := toml.Unmarshal([]byte(data), &decoded_content)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(
			0,
			fmt.Sprintf("The TOML file content cannot be decoded.\n\nOriginal Error: %s", err),
		)
		return
	}

	_, terraformValue, diags := convertToTerraformType(decoded_content)

	if diags.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Error = resp.Result.Set(ctx, types.DynamicValue(terraformValue))
}

func convertToTerraformType(dynamicValue any) (attr.Type, attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	switch value := dynamicValue.(type) {
	case string:
		return types.StringType, types.StringValue(value), diags
	case int:
		return types.Int64Type, types.Int64Value(int64(value)), diags
	case int64:
		return types.Int64Type, types.Int64Value(value), diags
	case float32:
		return types.Float64Type, types.Float64Value(float64(value)), diags
	case float64:
		return types.Float64Type, types.Float64Value(value), diags
	case bool:
		return types.BoolType, types.BoolValue(value), diags
	case time.Time:
		return types.StringType, types.StringValue(value.String()), diags
	case []any:
		elementTypes := make([]attr.Type, len(value))
		elementValues := make([]attr.Value, len(value))
		for i, dynamicElementValue := range value {
			elementType, elementValue, elementDiags := convertToTerraformType(dynamicElementValue)
			elementTypes[i] = elementType
			elementValues[i] = elementValue
			diags.Append(elementDiags...)
		}
		result, tupleDiags := types.TupleValue(elementTypes, elementValues)
		diags.Append(tupleDiags...)
		return types.TupleType{ElemTypes: elementTypes}, result, diags
	case map[string]any:
		attributeTypes := make(map[string]attr.Type, len(value))
		attributeValues := make(map[string]attr.Value, len(value))
		for attributeName, dynamicAttributeValue := range value {
			attributeType, attributeValue, attributeDiags := convertToTerraformType(dynamicAttributeValue)
			attributeTypes[attributeName] = attributeType
			attributeValues[attributeName] = attributeValue
			diags.Append(attributeDiags...)
		}
		result, objectDiags := types.ObjectValue(attributeTypes, attributeValues)
		diags.Append(objectDiags...)
		return types.ObjectType{AttrTypes: attributeTypes}, result, diags
	default:
		diags.AddError(
			"Invalid type to convert to Terraform type",
			fmt.Sprintf("Unable to convert value %v (type %T) to Terraform type", value, value),
		)
		return nil, nil, diags
	}
}
