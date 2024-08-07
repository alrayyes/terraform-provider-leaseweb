package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state model.Instance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Read public cloud instance %q",
		state.Id.ValueString(),
	))
	instance, err := i.client.PublicCloudFacade.GetInstance(
		state.Id.ValueString(),
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Instance", err.Error())

		logging.FacadeError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf("Unable to read instance %q", state.Id.ValueString()),
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
