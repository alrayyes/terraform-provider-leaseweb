package publiccloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &loadBalancerListenerResource{}
	_ resource.ResourceWithImportState = &loadBalancerListenerResource{}
)

type loadBalancerListenerDefaultRuleResourceModel struct {
	TargetGroupID types.String `tfsdk:"target_group_id"`
}

func (l loadBalancerListenerDefaultRuleResourceModel) generateLoadBalancerListenerDefaultRule() publicCloud.LoadBalancerListenerDefaultRule {
	return *publicCloud.NewLoadBalancerListenerDefaultRule(l.TargetGroupID.ValueString())
}

func (l loadBalancerListenerDefaultRuleResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_group_id": types.StringType,
	}
}

func adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource(sdkLoadBalancerListenerRule publicCloud.LoadBalancerListenerRule) loadBalancerListenerDefaultRuleResourceModel {
	return loadBalancerListenerDefaultRuleResourceModel{
		TargetGroupID: basetypes.NewStringValue(sdkLoadBalancerListenerRule.GetTargetGroupId()),
	}
}

type loadBalancerListenerCertificateResourceModel struct {
	PrivateKey  types.String `tfsdk:"private_key"`
	Certificate types.String `tfsdk:"certificate"`
	Chain       types.String `tfsdk:"chain"`
}

func (l loadBalancerListenerCertificateResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"private_key": types.StringType,
		"certificate": types.StringType,
		"chain":       types.StringType,
	}
}

func (l loadBalancerListenerCertificateResourceModel) generateSslCertificate() publicCloud.SslCertificate {
	sslCertificate := publicCloud.NewSslCertificate(
		l.PrivateKey.ValueString(),
		l.Certificate.ValueString(),
	)
	if !l.Chain.IsNull() && l.Chain.ValueString() != "" {
		sslCertificate.SetChain(l.Chain.ValueString())
	}

	return *sslCertificate
}

func adaptSslCertificateToLoadBalancerListenerCertificateResource(sdkSslCertificate publicCloud.SslCertificate) loadBalancerListenerCertificateResourceModel {
	listener := loadBalancerListenerCertificateResourceModel{
		PrivateKey:  basetypes.NewStringValue(sdkSslCertificate.GetPrivateKey()),
		Certificate: basetypes.NewStringValue(sdkSslCertificate.GetCertificate()),
	}

	chain, _ := sdkSslCertificate.GetChainOk()
	if chain != nil && *chain != "" {
		listener.Chain = basetypes.NewStringPointerValue(chain)
	}

	return listener
}

type LoadBalancerListenerResourceModel struct {
	ListenerID     types.String `tfsdk:"listener_id"`
	LoadBalancerID types.String `tfsdk:"load_balancer_id"`
	Protocol       types.String `tfsdk:"protocol"`
	Port           types.Int32  `tfsdk:"port"`
	Certificate    types.Object `tfsdk:"certificate"`
	DefaultRule    types.Object `tfsdk:"default_rule"`
}

func (l LoadBalancerListenerResourceModel) generateLoadBalancerListenerCreateOpts(ctx context.Context) (
	*publicCloud.LoadBalancerListenerCreateOpts,
	error,
) {
	defaultRule := loadBalancerListenerDefaultRuleResourceModel{}
	defaultRuleDiags := l.DefaultRule.As(ctx, &defaultRule, basetypes.ObjectAsOptions{})
	if defaultRuleDiags != nil {
		return nil, utils.ReturnError("generateLoadBalancerListenerCreateOpts", defaultRuleDiags)
	}

	opts := publicCloud.NewLoadBalancerListenerCreateOpts(
		publicCloud.Protocol(l.Protocol.ValueString()),
		l.Port.ValueInt32(),
		defaultRule.generateLoadBalancerListenerDefaultRule(),
	)

	if !l.Certificate.IsNull() {
		certificate := loadBalancerListenerCertificateResourceModel{}
		certificateDiags := l.Certificate.As(ctx, &certificate, basetypes.ObjectAsOptions{})
		if certificateDiags != nil {
			return nil, utils.ReturnError("generateLoadBalancerListenerCreateOpts", certificateDiags)
		}

		opts.SetCertificate(certificate.generateSslCertificate())
	}

	return opts, nil
}

func (l LoadBalancerListenerResourceModel) generateLoadBalancerListenerUpdateOpts(ctx context.Context) (
	*publicCloud.LoadBalancerListenerOpts,
	error,
) {
	opts := publicCloud.NewLoadBalancerListenerOpts()
	opts.SetProtocol(publicCloud.Protocol(l.Protocol.ValueString()))
	opts.SetPort(l.Port.ValueInt32())

	if !l.Certificate.IsNull() {
		certificate := loadBalancerListenerCertificateResourceModel{}
		certificateDiags := l.Certificate.As(
			ctx,
			&certificate,
			basetypes.ObjectAsOptions{},
		)
		if certificateDiags != nil {
			return nil, utils.ReturnError(
				"generateLoadBalancerListenerUpdateOpts",
				certificateDiags,
			)
		}

		opts.SetCertificate(certificate.generateSslCertificate())
	}

	if !l.DefaultRule.IsNull() {
		defaultRule := loadBalancerListenerDefaultRuleResourceModel{}
		defaultRuleDiags := l.DefaultRule.As(
			ctx,
			&defaultRule,
			basetypes.ObjectAsOptions{},
		)
		if defaultRuleDiags != nil {
			return nil, utils.ReturnError(
				"generateLoadBalancerListenerUpdateOpts",
				defaultRuleDiags,
			)
		}

		opts.SetDefaultRule(defaultRule.generateLoadBalancerListenerDefaultRule())
	}

	return opts, nil
}

func adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(
	sdkLoadBalancerListenerDetails publicCloud.LoadBalancerListenerDetails,
	ctx context.Context,
) (*LoadBalancerListenerResourceModel, error) {
	listener := LoadBalancerListenerResourceModel{
		ListenerID: basetypes.NewStringValue(sdkLoadBalancerListenerDetails.GetId()),
		Protocol:   basetypes.NewStringValue(string(sdkLoadBalancerListenerDetails.GetProtocol())),
		Port:       basetypes.NewInt32Value(sdkLoadBalancerListenerDetails.GetPort()),
	}

	if len(sdkLoadBalancerListenerDetails.SslCertificates) > 0 {
		certificate, err := utils.AdaptSdkModelToResourceObject(
			sdkLoadBalancerListenerDetails.SslCertificates[0],
			loadBalancerListenerCertificateResourceModel{}.attributeTypes(),
			ctx,
			adaptSslCertificateToLoadBalancerListenerCertificateResource,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource: %w",
				err,
			)
		}
		listener.Certificate = certificate
	}

	if len(sdkLoadBalancerListenerDetails.Rules) > 0 {
		defaultRule, err := utils.AdaptSdkModelToResourceObject(
			sdkLoadBalancerListenerDetails.Rules[0],
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			ctx,
			adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource: %w",
				err,
			)
		}
		listener.DefaultRule = defaultRule
	}

	return &listener, nil
}

func adaptLoadBalancerListenerToLoadBalancerListenerResource(
	sdkLoadBalancerListener publicCloud.LoadBalancerListener,
	ctx context.Context,
) (*LoadBalancerListenerResourceModel, error) {
	listener := LoadBalancerListenerResourceModel{
		ListenerID: basetypes.NewStringValue(sdkLoadBalancerListener.GetId()),
		Protocol:   basetypes.NewStringValue(string(sdkLoadBalancerListener.Protocol)),
		Port:       basetypes.NewInt32Value(sdkLoadBalancerListener.GetPort()),
	}

	if len(sdkLoadBalancerListener.Rules) > 0 {
		defaultRule, err := utils.AdaptSdkModelToResourceObject(
			sdkLoadBalancerListener.Rules[0],
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			ctx,
			adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"adaptLoadBalancerListenerToLoadBalancerListenerResource: %w",
				err,
			)
		}
		listener.DefaultRule = defaultRule
	}

	return &listener, nil
}

type loadBalancerListenerResource struct {
	name   string
	client publicCloud.PublicCloudAPI
}

func (l *loadBalancerListenerResource) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
	if request.ProviderData == nil {
		return
	}

	coreClient, ok := request.ProviderData.(client.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				request.ProviderData,
			),
		)

		return
	}

	l.client = coreClient.PublicCloudAPI
}

func (l *loadBalancerListenerResource) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, l.name)
}

func (l *loadBalancerListenerResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer ID",
			},
			"listener_id": schema.StringAttribute{
				Computed:    true,
				Description: "Listener ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Valid options are " + utils.StringTypeArrayToMarkdown(publicCloud.AllowedProtocolEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedProtocolEnumValues)...),
				},
			},
			"port": schema.Int32Attribute{
				Required:    true,
				Description: "Port that the listener listens to",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"certificate": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Required only if protocol is HTTPS",
				Attributes: map[string]schema.Attribute{
					"private_key": schema.StringAttribute{
						Optional:    true,
						Description: "Client Private Key. Required only if protocol is `HTTPS`",
						Sensitive:   true,
					},
					"certificate": schema.StringAttribute{
						Optional:    true,
						Description: "Client Certificate. Required only if protocol is `HTTPS`",
						Sensitive:   true,
					},
					"chain": schema.StringAttribute{
						Optional:    true,
						Description: "CA certificate. Not required, but can be added if protocol is `HTTPS`",
						Sensitive:   true,
					},
				},
			},
			"default_rule": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"target_group_id": schema.StringAttribute{
						Optional:    true,
						Description: "Client Private Key. Required only if protocol is `HTTPS`",
					},
				},
			},
		},
	}
}

func (l *loadBalancerListenerResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan LoadBalancerListenerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	summary := fmt.Sprintf("Creating resource %s", l.name)

	opts, err := plan.generateLoadBalancerListenerCreateOpts(ctx)
	if err != nil {
		response.Diagnostics.AddError(summary, utils.DefaultErrMsg)
		return
	}

	sdkLoadBalancerListener, httpResponse, err := l.client.CreateLoadBalancerListener(
		ctx,
		plan.LoadBalancerID.ValueString(),
	).LoadBalancerListenerCreateOpts(*opts).Execute()
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	state, resourceErr := adaptLoadBalancerListenerToLoadBalancerListenerResource(
		*sdkLoadBalancerListener,
		ctx,
	)
	if resourceErr != nil {
		response.Diagnostics.AddError(summary, utils.DefaultErrMsg)
		return
	}

	state.LoadBalancerID = plan.LoadBalancerID
	state.Certificate = plan.Certificate

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerListenerResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var requestState LoadBalancerListenerResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &requestState)...)
	if response.Diagnostics.HasError() {
		return
	}

	summary := fmt.Sprintf(
		"Reading resource %s for load_balancer_id %q listener_id %q",
		l.name,
		requestState.LoadBalancerID.ValueString(),
		requestState.ListenerID.ValueString(),
	)

	sdkLoadBalancerListenerDetails, httpResponse, err := l.client.GetLoadBalancerListener(
		ctx,
		requestState.LoadBalancerID.ValueString(),
		requestState.ListenerID.ValueString(),
	).Execute()
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	state, resourceErr := adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(*sdkLoadBalancerListenerDetails, ctx)
	if resourceErr != nil {
		utils.Error(ctx, &response.Diagnostics, summary, resourceErr, nil)
		return
	}

	state.LoadBalancerID = requestState.LoadBalancerID
	state.ListenerID = requestState.ListenerID
	if state.Certificate.IsNull() {
		state.Certificate = requestState.Certificate
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerListenerResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan LoadBalancerListenerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	summary := fmt.Sprintf(
		"Updating resource %s for load_balancer_id %q listener_id %q",
		l.name,
		plan.LoadBalancerID.ValueString(),
		plan.ListenerID.ValueString(),
	)

	opts, err := plan.generateLoadBalancerListenerUpdateOpts(ctx)
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, nil)
		return
	}

	sdkLoadBalancerListener, httpResponse, err := l.client.
		UpdateLoadBalancerListener(
			ctx,
			plan.LoadBalancerID.ValueString(),
			plan.ListenerID.ValueString(),
		).
		LoadBalancerListenerOpts(*opts).
		Execute()
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	state, resourceErr := adaptLoadBalancerListenerToLoadBalancerListenerResource(
		*sdkLoadBalancerListener,
		ctx,
	)
	if resourceErr != nil {
		utils.Error(ctx, &response.Diagnostics, summary, resourceErr, nil)
		return
	}

	if state.Certificate.IsNull() {
		state.Certificate = plan.Certificate
	}
	state.LoadBalancerID = plan.LoadBalancerID

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerListenerResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state LoadBalancerListenerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := l.client.DeleteLoadBalancerListener(
		ctx,
		state.LoadBalancerID.ValueString(),
		state.ListenerID.ValueString(),
	).Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Deleting resource %s for instance_id %q listener_id %q",
			l.name,
			state.LoadBalancerID.ValueString(),
			state.ListenerID.ValueString(),
		)
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
	}
}

func (l *loadBalancerListenerResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf(
				"Expected import identifier with format: load_balancer_id,listener_id. Got: %q",
				request.ID,
			),
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("load_balancer_id"),
		idParts[0],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("listener_id"),
		idParts[1],
	)...)
}

func NewLoadBalancerListenerResource() resource.Resource {
	return &loadBalancerListenerResource{
		name: "public_cloud_load_balancer_listener",
	}
}