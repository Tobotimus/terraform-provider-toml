// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure tomlProvider satisfies various provider interfaces.
var _ provider.Provider = &tomlProvider{}

// tomlProvider defines the provider implementation.
type tomlProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *tomlProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "toml"
	resp.Version = p.version
}

func (p *tomlProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {

}

func (p *tomlProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *tomlProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *tomlProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTomlFileDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &tomlProvider{
			version: version,
		}
	}
}
