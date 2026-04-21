package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &FlagDataSource{}

type FlagDataSource struct {
	client *flagdClient
}

type FlagDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewFlagDataSource() datasource.DataSource {
	return &FlagDataSource{}
}

func (d *FlagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flag"
}

func (d *FlagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a feature flag from flagd. Use this to reference an existing flag without managing it.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the flag to look up.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the flag is currently enabled.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the flag.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the flag was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the flag was last updated.",
				Computed:    true,
			},
		},
	}
}

func (d *FlagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*flagdClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type",
			fmt.Sprintf("Expected *flagdClient, got %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *FlagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state FlagDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag, err := d.client.GetFlag(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading flag", err.Error())
		return
	}
	if flag == nil {
		resp.Diagnostics.AddError("Flag not found",
			fmt.Sprintf("Flag %q does not exist on the flagd server.", state.Name.ValueString()))
		return
	}

	state.Enabled = types.BoolValue(flag.Enabled)
	state.Description = types.StringValue(flag.Description)
	state.CreatedAt = types.StringValue(flag.CreatedAt.String())
	state.UpdatedAt = types.StringValue(flag.UpdatedAt.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
