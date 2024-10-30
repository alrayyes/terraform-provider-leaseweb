package dedicatedserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &bandwidthNotificationSettingResource{}
	_ resource.ResourceWithConfigure = &bandwidthNotificationSettingResource{}
)

type bandwidthNotificationSettingResource struct {
	client dedicatedServer.DedicatedServerAPI
}

type bandwidthNotificationSettingResourceModel struct {
	Id                types.String `tfsdk:"id"`
	DedicatedServerId types.String `tfsdk:"dedicated_server_id"`
	Frequency         types.String `tfsdk:"frequency"`
	Threshold         types.String `tfsdk:"threshold"`
	Unit              types.String `tfsdk:"unit"`
}

func NewBandwidthNotificationSettingResource() resource.Resource {
	return &bandwidthNotificationSettingResource{}
}

func (b *bandwidthNotificationSettingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_bandwidth_notification_setting"
}

func (b *bandwidthNotificationSettingResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	coreClient, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	b.client = coreClient.DedicatedServerAPI
}

func (b *bandwidthNotificationSettingResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The notification setting bandwidth unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dedicated_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The server unique identifier",
			},
			"frequency": schema.StringAttribute{
				Required:    true,
				Description: "The notification frequency. Valid options can be *DAILY* or *WEEKLY* or *MONTHLY*.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DAILY", "WEEKLY", "MONTHLY"}...),
				},
			},
			"threshold": schema.StringAttribute{
				Required:    true,
				Description: "Threshold Value. Value can be a number greater than 0.",
				Validators: []validator.String{
					greaterThanZero(),
				},
			},
			"unit": schema.StringAttribute{
				Required:    true,
				Description: "The notification unit. Valid options can be *Mbps* or *Gbps*.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Mbps", "Gbps"}...),
				},
			},
		},
	}
}

func (b *bandwidthNotificationSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data bandwidthNotificationSettingResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewBandwidthNotificationSettingOpts(
		data.Frequency.ValueString(),
		data.Threshold.ValueString(),
		data.Unit.ValueString(),
	)
	request := b.client.CreateServerBandwidthNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
	).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error creating bandwidth notification setting with dedicated_server_id: %q", data.DedicatedServerId.ValueString())
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	newData := bandwidthNotificationSettingResourceModel{
		Id:        types.StringValue(result.GetId()),
		Frequency: types.StringValue(result.GetFrequency()),
		Threshold: types.StringValue(result.GetThreshold()),
		Unit:      types.StringValue(result.GetUnit()),
	}
	newData.DedicatedServerId = data.DedicatedServerId
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (b *bandwidthNotificationSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data bandwidthNotificationSettingResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := b.client.GetServerBandwidthNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
		data.Id.ValueString(),
	)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error reading bandwidth notification setting with id: %q and dedicated_server_id: %q", data.Id.ValueString(), data.DedicatedServerId.ValueString())
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	newData := bandwidthNotificationSettingResourceModel{
		Id:        types.StringValue(result.GetId()),
		Frequency: types.StringValue(result.GetFrequency()),
		Threshold: types.StringValue(result.GetThreshold()),
		Unit:      types.StringValue(result.GetUnit()),
	}
	newData.DedicatedServerId = data.DedicatedServerId
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (b *bandwidthNotificationSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bandwidthNotificationSettingResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewBandwidthNotificationSettingOpts(
		data.Frequency.ValueString(),
		data.Threshold.ValueString(),
		data.Unit.ValueString(),
	)
	request := b.client.UpdateServerBandwidthNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
		data.Id.ValueString(),
	).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error updating bandwidth notification setting with id: %q and dedicated_server_id: %q", data.Id.ValueString(), data.DedicatedServerId.ValueString())
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	newData := bandwidthNotificationSettingResourceModel{
		Id:                data.Id,
		DedicatedServerId: data.DedicatedServerId,
		Frequency:         types.StringValue(result.GetFrequency()),
		Threshold:         types.StringValue(result.GetThreshold()),
		Unit:              types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (b *bandwidthNotificationSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bandwidthNotificationSettingResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := b.client.DeleteServerBandwidthNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
		data.Id.ValueString(),
	)
	response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error deleting bandwidth notification setting with id: %q and dedicated_server_id: %q", data.Id.ValueString(), data.DedicatedServerId.ValueString())
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}
}
