// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure tomlProvider satisfies various provider interfaces.
var _ provider.Provider = &TomlProvider{}

// TomlProvider defines the provider implementation.
type TomlProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *TomlProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "toml"
	resp.Version = p.version
}

func (p *TomlProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {

}

func (p *TomlProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *TomlProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *TomlProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTomlFileDataSource,
	}
}

func (p *TomlProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewDecodeFunction,
		NewEncodeFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TomlProvider{
			version: version,
		}
	}
}
