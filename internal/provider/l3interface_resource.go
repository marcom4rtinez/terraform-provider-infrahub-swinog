package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	infrahub_sdk "github.com/opsmill/infrahub-sdk-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &l3interfaceResource{}
	_ resource.ResourceWithConfigure = &l3interfaceResource{}
)

// NewL3interfaceResource is a helper function to simplify the provider implementation.
func NewL3interfaceResource() resource.Resource {
	return &l3interfaceResource{}
}

// l3interfaceResource is the resource implementation.
type l3interfaceResource struct {
	client                       *graphql.Client
	Edges_node_id                types.String `tfsdk:"id"`
	Edges_node_role_id           types.String `tfsdk:"role_id"`
	Edges_node_role_value        types.String `tfsdk:"role_value"`
	Edges_node_name_id           types.String `tfsdk:"name_id"`
	Edges_node_name_value        types.String `tfsdk:"name_value"`
	Edges_node_enabled_id        types.String `tfsdk:"enabled_id"`
	Edges_node_enabled_value     types.Bool   `tfsdk:"enabled_value"`
	Edges_node_description_id    types.String `tfsdk:"description_id"`
	Edges_node_description_value types.String `tfsdk:"description_value"`
	Edges_node_device_node_id    types.String `tfsdk:"device_node_id"`
	Edges_node_status_value      types.String `tfsdk:"status_value"`
	Edges_node_status_id         types.String `tfsdk:"status_id"`
}

// Metadata returns the resource type name.
func (r *l3interfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l3interface"
}

// Schema defines the schema for the resource.
func (r *l3interfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"role_id": schema.StringAttribute{
				Computed: true,
			},
			"name_id": schema.StringAttribute{
				Computed: true,
			},
			"enabled_id": schema.StringAttribute{
				Computed: true,
			},
			"description_id": schema.StringAttribute{
				Computed: true,
			},
			"status_id": schema.StringAttribute{
				Computed: true,
			},
			"role_value": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"name_value": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"enabled_value": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"description_value": schema.StringAttribute{
				Required: true,
			},
			"device_node_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"status_value": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *l3interfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan l3interfaceResource
	tflog.Info(ctx, req.Config.Raw.String())
	tflog.Info(ctx, req.Plan.Raw.String())
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var defaultL3interface infrahub_sdk.InfraInterfaceL3CreateInput

	// Assign each field, using the helper function to handle defaults
	defaultL3interface.Role.Value = plan.Edges_node_role_value.ValueString()
	defaultL3interface.Name.Value = plan.Edges_node_name_value.ValueString()
	defaultL3interface.Enabled.Value = plan.Edges_node_enabled_value.ValueBool()
	defaultL3interface.Description.Value = plan.Edges_node_description_value.ValueString()
	defaultL3interface.Device.Id = plan.Edges_node_device_node_id.ValueString()
	defaultL3interface.Status.Value = plan.Edges_node_status_value.ValueString()

	tflog.Info(ctx, fmt.Sprint("Creating L3interface ", plan.Edges_node_description_value))

	response, err := infrahub_sdk.L3interfaceCreate(ctx, *r.client, defaultL3interface)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create l3interface in Infrahub",
			err.Error(),
		)
		return
	}
	plan.Edges_node_id = types.StringValue(response.InfraInterfaceL3Create.Object.GetId())
	plan.Edges_node_role_id = types.StringValue(response.InfraInterfaceL3Create.Object.Role.GetId())
	plan.Edges_node_role_value = types.StringValue(response.InfraInterfaceL3Create.Object.Role.Value)
	plan.Edges_node_name_id = types.StringValue(response.InfraInterfaceL3Create.Object.Name.GetId())
	plan.Edges_node_name_value = types.StringValue(response.InfraInterfaceL3Create.Object.Name.Value)
	plan.Edges_node_enabled_id = types.StringValue(response.InfraInterfaceL3Create.Object.Enabled.GetId())
	plan.Edges_node_enabled_value = types.BoolValue(response.InfraInterfaceL3Create.Object.Enabled.Value)
	plan.Edges_node_description_id = types.StringValue(response.InfraInterfaceL3Create.Object.Description.GetId())
	plan.Edges_node_description_value = types.StringValue(response.InfraInterfaceL3Create.Object.Description.Value)
	plan.Edges_node_device_node_id = types.StringValue(response.InfraInterfaceL3Create.Object.Device.Node.GetId())
	plan.Edges_node_status_value = types.StringValue(response.InfraInterfaceL3Create.Object.Status.Value)
	plan.Edges_node_status_id = types.StringValue(response.InfraInterfaceL3Create.Object.Status.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *l3interfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading L3interface...")
	var state l3interfaceResource

	// Read configuration into config
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprint("Reading L3interface ", state.Edges_node_description_value))

	// Call the API with the specified device_name from the configuration
	response, err := infrahub_sdk.L3interface(ctx, *r.client, []string{state.Edges_node_description_value.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read l3interface from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.InfraInterfaceL3.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single l3interface, query didn't return exactly 1 l3interface",
			"Expected exactly 1 l3interface in response, got a different count.",
		)
		return
	}
	state.Edges_node_id = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.GetId())
	state.Edges_node_role_id = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Role.GetId())
	state.Edges_node_role_value = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Role.Value)
	state.Edges_node_name_id = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Name.GetId())
	state.Edges_node_name_value = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Name.Value)
	state.Edges_node_enabled_id = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Enabled.GetId())
	state.Edges_node_enabled_value = types.BoolValue(response.InfraInterfaceL3.Edges[0].Node.Enabled.Value)
	state.Edges_node_description_id = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Description.GetId())
	state.Edges_node_description_value = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Description.Value)
	state.Edges_node_device_node_id = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Device.Node.GetId())
	state.Edges_node_status_value = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Status.Value)
	state.Edges_node_status_id = types.StringValue(response.InfraInterfaceL3.Edges[0].Node.Status.GetId())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *l3interfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve the planned configuration values from Terraform
	var plan l3interfaceResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current state
	var state l3interfaceResource
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateInput infrahub_sdk.InfraInterfaceL3UpsertInput

	// Prepare the update input using values from the plan and applying defaults
	updateInput.Role.Value = setDefault(plan.Edges_node_role_value.ValueString(), state.Edges_node_role_value.ValueString())
	updateInput.Name.Value = setDefault(plan.Edges_node_name_value.ValueString(), state.Edges_node_name_value.ValueString())
	updateInput.Enabled.Value = plan.Edges_node_enabled_value.ValueBool()
	updateInput.Description.Value = setDefault(plan.Edges_node_description_value.ValueString(), state.Edges_node_description_value.ValueString())
	updateInput.Device.Id = setDefault(plan.Edges_node_device_node_id.ValueString(), state.Edges_node_device_node_id.ValueString())
	updateInput.Status.Value = setDefault(plan.Edges_node_status_value.ValueString(), state.Edges_node_status_value.ValueString())
	updateInput.Id = state.Edges_node_id.ValueString()

	// Log the update operation
	tflog.Info(ctx, fmt.Sprintf("Updating L3interface %s", state.Edges_node_description_value.ValueString()))

	// Send the update request to the API
	response, err := infrahub_sdk.L3interfaceUpsert(ctx, *r.client, updateInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update device in Infrahub",
			err.Error(),
		)
		return
	}
	plan.Edges_node_id = types.StringValue(response.InfraInterfaceL3Upsert.Object.GetId())
	plan.Edges_node_role_id = types.StringValue(response.InfraInterfaceL3Upsert.Object.Role.GetId())
	plan.Edges_node_role_value = types.StringValue(response.InfraInterfaceL3Upsert.Object.Role.Value)
	plan.Edges_node_name_id = types.StringValue(response.InfraInterfaceL3Upsert.Object.Name.GetId())
	plan.Edges_node_name_value = types.StringValue(response.InfraInterfaceL3Upsert.Object.Name.Value)
	plan.Edges_node_enabled_id = types.StringValue(response.InfraInterfaceL3Upsert.Object.Enabled.GetId())
	plan.Edges_node_enabled_value = types.BoolValue(response.InfraInterfaceL3Upsert.Object.Enabled.Value)
	plan.Edges_node_description_id = types.StringValue(response.InfraInterfaceL3Upsert.Object.Description.GetId())
	plan.Edges_node_description_value = types.StringValue(response.InfraInterfaceL3Upsert.Object.Description.Value)
	plan.Edges_node_device_node_id = types.StringValue(response.InfraInterfaceL3Upsert.Object.Device.Node.GetId())
	plan.Edges_node_status_value = types.StringValue(response.InfraInterfaceL3Upsert.Object.Status.Value)
	plan.Edges_node_status_id = types.StringValue(response.InfraInterfaceL3Upsert.Object.Status.GetId())

	// Set the updated state with the latest data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *l3interfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state l3interfaceResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := infrahub_sdk.L3interfaceDelete(ctx, *r.client, state.Edges_node_id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting L3interface",
			"Could not delete l3interface, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *l3interfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &client
}
