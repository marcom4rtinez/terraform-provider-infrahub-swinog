package provider

import (
	"context"
	"fmt"

	infrahub_sdk "github.com/opsmill/infrahub-sdk-go"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ipaddressqueryDataSource{}
	_ datasource.DataSourceWithConfigure = &ipaddressqueryDataSource{}
)

// NewIpaddressqueryDataSource is a helper function to simplify the provider implementation.
func NewIpaddressqueryDataSource() datasource.DataSource {
	return &ipaddressqueryDataSource{}
}

type ipaddressqueryDataSource struct {
	client                           *graphql.Client
	Ip_address_value                 types.String `tfsdk:"ip_address_value"`
	Edges_node_id                    types.String `tfsdk:"id"`
	Edges_node_address_value         types.String `tfsdk:"address_value"`
	Edges_node_address_ip            types.String `tfsdk:"address_ip"`
	Edges_node_address_netmask       types.String `tfsdk:"address_netmask"`
	Edges_node_address_with_hostmask types.String `tfsdk:"address_with_hostmask"`
	Edges_node_address_with_netmask  types.String `tfsdk:"address_with_netmask"`
	Edges_node_description_value     types.String `tfsdk:"description_value"`
}

func (d *ipaddressqueryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipaddressquery"
}

func (d *ipaddressqueryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ip_address_value": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"address_value": schema.StringAttribute{
				Computed: true,
			},
			"address_ip": schema.StringAttribute{
				Computed: true,
			},
			"address_netmask": schema.StringAttribute{
				Computed: true,
			},
			"address_with_hostmask": schema.StringAttribute{
				Computed: true,
			},
			"address_with_netmask": schema.StringAttribute{
				Computed: true,
			},
			"description_value": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *ipaddressqueryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading ipaddressquery data...")
	var config ipaddressqueryDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Ipaddressquery(ctx, *d.client, config.Ip_address_value.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ipaddressquery from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.InfraIPAddress.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single ipaddressquery, query didn't return exactly 1 ipaddressquery",
			"Expected exactly 1 ipaddressquery in response, got a different count.",
		)
		return
	}

	state := ipaddressqueryDataSource{
		Ip_address_value:                 config.Ip_address_value,
		Edges_node_id:                    types.StringValue(response.InfraIPAddress.Edges[0].Node.Id),
		Edges_node_address_value:         types.StringValue(response.InfraIPAddress.Edges[0].Node.Address.Value),
		Edges_node_address_ip:            types.StringValue(response.InfraIPAddress.Edges[0].Node.Address.Ip),
		Edges_node_address_netmask:       types.StringValue(response.InfraIPAddress.Edges[0].Node.Address.Netmask),
		Edges_node_address_with_hostmask: types.StringValue(response.InfraIPAddress.Edges[0].Node.Address.With_hostmask),
		Edges_node_address_with_netmask:  types.StringValue(response.InfraIPAddress.Edges[0].Node.Address.With_netmask),
		Edges_node_description_value:     types.StringValue(response.InfraIPAddress.Edges[0].Node.Description.Value),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *ipaddressqueryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = &client
}
