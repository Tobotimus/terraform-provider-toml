package provider

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pelletier/go-toml/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &tomlfileDataSource{}
)

// NewTomlFileDataSource is a helper function to simplify the provider implementation.
func NewTomlFileDataSource() datasource.DataSource {
	return &tomlfileDataSource{}
}

// tomlfileDataSource is the data source implementation.
type tomlfileDataSource struct{}

// Metadata returns the data source type name.
func (d *tomlfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

// Schema defines the schema for the data source.
func (d *tomlfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The `toml_file` data source allows Terraform to parse TOML file content as a data source.\n\n" +
			"Note that the content is returned in `content_json` as a JSON-encoded string. You can then use " +
			"the `jsondecode()` function to access fields within the TOML file. This oddity is due to a current " +
			"limitation in the Terraform Plugin Framework.",
		Attributes: map[string]schema.Attribute{
			"input": schema.StringAttribute{
				Description: "Raw content of the TOML file to be parsed.",
				Required:    true,
			},
			"content": schema.DynamicAttribute{
				Description: "Decoded content of the TOML file.",
				Computed:    true,
			},
			"content_json": schema.StringAttribute{
				Description: "JSON-encoded content of the TOML file.",
				Computed:    true,
				DeprecationMessage: "The `content_json` attribute is deprecated, and will be removed in the next " +
					"major version. Use the `content` attribute instead.",
			},
			"id": schema.StringAttribute{
				Description: "The hexadecimal encoding of the SHA1 checksum of the JSON-encoded content.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *tomlfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config tomlFileDataSourceModelV0

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var decoded_content interface{}
	err := toml.Unmarshal([]byte(config.Input.ValueString()), &decoded_content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read TOML file data source error",
			"The TOML file content cannot be decoded.\n\n"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	content_json, err := json.Marshal(decoded_content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read TOML file data source error",
			"The loaded content of the TOML file could not be encoded as JSON.\n\n"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	_, content_tf, diags := convertToTerraformType(decoded_content)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sha1Sum := sha1.Sum(content_json)
	sha1Hex := hex.EncodeToString(sha1Sum[:])

	state := tomlFileDataSourceModelV0{
		Input:       config.Input,
		Content:     types.DynamicValue(content_tf),
		ContentJSON: types.StringValue(string(content_json)),
		ID:          types.StringValue(sha1Hex),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

type tomlFileDataSourceModelV0 struct {
	Input       types.String  `tfsdk:"input"`
	Content     types.Dynamic `tfsdk:"content"`
	ContentJSON types.String  `tfsdk:"content_json"`
	ID          types.String  `tfsdk:"id"`
}
