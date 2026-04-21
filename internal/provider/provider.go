package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &FlagdProvider{}

type FlagdProvider struct {
	version string
}

type FlagdProviderModel struct {
	URL      types.String `tfsdk:"url"`
	AdminKey types.String `tfsdk:"admin_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FlagdProvider{version: version}
	}
}

func (p *FlagdProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "flagd"
	resp.Version = p.version
}

func (p *FlagdProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage feature flags on a flagd server.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "The flagd server URL. Can also be set via FLAGD_URL environment variable.",
				Optional:    true,
			},
			"admin_key": schema.StringAttribute{
				Description: "Admin key for write access. Can also be set via FLAGD_ADMIN_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *FlagdProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config FlagdProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("FLAGD_URL")
	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}
	if url == "" {
		url = "http://localhost:8080"
	}

	adminKey := os.Getenv("FLAGD_ADMIN_KEY")
	if !config.AdminKey.IsNull() {
		adminKey = config.AdminKey.ValueString()
	}

	client := newFlagdClient(url, adminKey)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *FlagdProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFlagResource,
	}
}

func (p *FlagdProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFlagDataSource,
	}
}
