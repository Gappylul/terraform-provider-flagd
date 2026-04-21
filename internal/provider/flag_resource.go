package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &FlagResource{}

type FlagResource struct {
	client *flagdClient
}

type FlagResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewFlagResource() resource.Resource {
	return &FlagResource{}
}

func (r *FlagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flag"
}

func (r *FlagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a feature flag on a flagd server.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique name of the flag.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the flag is enabled",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: "Human-readable description of the flag.",
				Optional:    true,
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

func (r *FlagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*flagdClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type",
			fmt.Sprintf("Expected *flagdClient, got %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *FlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FlagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag, err := r.client.UpsertFlag(ctx,
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		plan.Enabled.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating flag", err.Error())
		return
	}

	plan.CreatedAt = types.StringValue(flag.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(flag.UpdatedAt.String())
	plan.Description = types.StringValue(flag.Description)
	plan.Enabled = types.BoolValue(flag.Enabled)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FlagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag, err := r.client.GetFlag(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading flag", err.Error())
		return
	}
	if flag == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Enabled = types.BoolValue(flag.Enabled)
	state.Description = types.StringValue(flag.Description)
	state.CreatedAt = types.StringValue(flag.CreatedAt.String())
	state.UpdatedAt = types.StringValue(flag.UpdatedAt.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FlagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag, err := r.client.UpsertFlag(ctx,
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		plan.Enabled.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating flag", err.Error())
		return
	}

	plan.CreatedAt = types.StringValue(flag.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(flag.UpdatedAt.String())
	plan.Description = types.StringValue(flag.Description)
	plan.Enabled = types.BoolValue(flag.Enabled)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FlagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteFlag(ctx, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting flag", err.Error())
	}
}
