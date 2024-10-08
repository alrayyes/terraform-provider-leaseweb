// Package to_data_source_model implements adapters to convert domain entities to sdk models.
package to_data_source_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
)

func AdaptInstances(domainInstances public_cloud.Instances) model.Instances {
	var instances model.Instances

	for _, domainInstance := range domainInstances {
		instance := adaptInstance(domainInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func adaptInstance(domainInstance public_cloud.Instance) model.Instance {
	instance := model.Instance{
		Id:     basetypes.NewStringValue(domainInstance.Id),
		Region: *adaptRegion(domainInstance.Region),
		Reference: shared.AdaptNullableStringToStringValue(
			domainInstance.Reference,
		),
		Resources: adaptResources(
			domainInstance.Resources,
		),
		Image:         adaptImage(domainInstance.Image),
		State:         basetypes.NewStringValue(string(domainInstance.State)),
		ProductType:   basetypes.NewStringValue(domainInstance.ProductType),
		HasPublicIpv4: basetypes.NewBoolValue(domainInstance.HasPublicIpv4),
		HasPrivateNetwork: basetypes.NewBoolValue(
			domainInstance.HasPrivateNetwork,
		),
		Type: adaptInstanceType(domainInstance.Type),
		RootDiskSize: basetypes.NewInt64Value(
			int64(domainInstance.RootDiskSize.Value),
		),
		RootDiskStorageType: basetypes.NewStringValue(
			string(domainInstance.RootDiskStorageType),
		),
		StartedAt: shared.AdaptNullableTimeToStringValue(
			domainInstance.StartedAt,
		),
		Contract: adaptContract(
			domainInstance.Contract,
		),
		MarketAppId: shared.AdaptNullableStringToStringValue(
			domainInstance.MarketAppId,
		),
		AutoScalingGroup: shared.AdaptNullableDomainEntityToDatasourceModel(
			domainInstance.AutoScalingGroup,
			adaptAutoScalingGroup,
		),
		Iso: shared.AdaptNullableDomainEntityToDatasourceModel(
			domainInstance.Iso,
			adaptIso,
		),
		PrivateNetwork: shared.AdaptNullableDomainEntityToDatasourceModel(
			domainInstance.PrivateNetwork,
			adaptPrivateNetwork,
		),
	}

	for _, autoScalingGroupIp := range domainInstance.Ips {
		ip := adaptIp(autoScalingGroupIp)
		instance.Ips = append(instance.Ips, ip)
	}

	return instance
}

func adaptResources(resources public_cloud.Resources) model.Resources {
	return model.Resources{
		Cpu:    adaptCpu(resources.Cpu),
		Memory: adaptMemory(resources.Memory),
		PublicNetworkSpeed: adaptNetworkSpeed(
			resources.PublicNetworkSpeed,
		),
		PrivateNetworkSpeed: adaptNetworkSpeed(
			resources.PrivateNetworkSpeed,
		),
	}
}

func adaptCpu(cpu public_cloud.Cpu) model.Cpu {
	return model.Cpu{
		Value: basetypes.NewInt64Value(int64(cpu.Value)),
		Unit:  basetypes.NewStringValue(cpu.Unit),
	}
}

func adaptMemory(memory public_cloud.Memory) model.Memory {
	return model.Memory{
		Value: basetypes.NewFloat64Value(memory.Value),
		Unit:  basetypes.NewStringValue(memory.Unit),
	}
}

func adaptNetworkSpeed(networkSpeed public_cloud.NetworkSpeed) model.NetworkSpeed {
	return model.NetworkSpeed{
		Value: basetypes.NewInt64Value(int64(networkSpeed.Value)),
		Unit:  basetypes.NewStringValue(networkSpeed.Unit),
	}
}

func adaptImage(domainImage public_cloud.Image) model.Image {
	image := model.Image{
		Id:           basetypes.NewStringValue(domainImage.Id),
		Name:         basetypes.NewStringValue(domainImage.Name),
		Version:      shared.AdaptNullableStringToStringValue(domainImage.Version),
		Family:       basetypes.NewStringValue(domainImage.Family),
		Flavour:      basetypes.NewStringValue(domainImage.Flavour),
		Architecture: shared.AdaptNullableStringToStringValue(domainImage.Architecture),
		State:        shared.AdaptNullableStringToStringValue(domainImage.State),
		StateReason:  shared.AdaptNullableStringToStringValue(domainImage.StateReason),
		Region: shared.AdaptNullableDomainEntityToDatasourceModel(
			domainImage.Region,
			adaptRegion,
		),
		CreatedAt: shared.AdaptNullableTimeToStringValue(domainImage.CreatedAt),
		UpdatedAt: shared.AdaptNullableTimeToStringValue(domainImage.UpdatedAt),
		Custom:    shared.AdaptBoolToBoolValue(domainImage.Custom),
		StorageSize: shared.AdaptNullableDomainEntityToDatasourceModel(
			domainImage.StorageSize,
			adaptStorageSize,
		),
	}

	for _, marketApp := range domainImage.MarketApps {
		image.MarketApps = append(
			image.MarketApps, types.StringValue(marketApp),
		)
	}

	for _, storageType := range domainImage.StorageTypes {
		image.StorageTypes = append(
			image.StorageTypes, types.StringValue(storageType),
		)
	}

	return image
}

func adaptContract(contract public_cloud.Contract) model.Contract {
	return model.Contract{
		BillingFrequency: basetypes.NewInt64Value(
			int64(contract.BillingFrequency),
		),
		Term:       basetypes.NewInt64Value(int64(contract.Term)),
		Type:       basetypes.NewStringValue(string(contract.Type)),
		EndsAt:     shared.AdaptNullableTimeToStringValue(contract.EndsAt),
		RenewalsAt: basetypes.NewStringValue(contract.RenewalsAt.String()),
		CreatedAt:  basetypes.NewStringValue(contract.CreatedAt.String()),
		State:      basetypes.NewStringValue(string(contract.State)),
	}
}

func adaptAutoScalingGroup(autoScalingGroup public_cloud.AutoScalingGroup) *model.AutoScalingGroup {
	return &model.AutoScalingGroup{
		Id:    basetypes.NewStringValue(autoScalingGroup.Id),
		Type:  basetypes.NewStringValue(string(autoScalingGroup.Type)),
		State: basetypes.NewStringValue(string(autoScalingGroup.State)),
		DesiredAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.DesiredAmount,
		),
		Region: *adaptRegion(autoScalingGroup.Region),
		Reference: basetypes.NewStringValue(
			autoScalingGroup.Reference.String(),
		),
		CreatedAt: basetypes.NewStringValue(
			autoScalingGroup.CreatedAt.String(),
		),
		UpdatedAt: basetypes.NewStringValue(
			autoScalingGroup.UpdatedAt.String(),
		),
		StartsAt: shared.AdaptNullableTimeToStringValue(
			autoScalingGroup.StartsAt,
		),
		EndsAt: shared.AdaptNullableTimeToStringValue(
			autoScalingGroup.EndsAt,
		),
		MinimumAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.MinimumAmount,
		),
		MaximumAmount: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.MaximumAmount,
		),
		CpuThreshold: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.CpuThreshold,
		),
		WarmupTime: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.WarmupTime,
		),
		CooldownTime: shared.AdaptNullableIntToInt64Value(
			autoScalingGroup.CooldownTime,
		),
	}
}

func adaptLoadBalancer(loadBalancer public_cloud.LoadBalancer) *model.LoadBalancer {
	var ips []model.Ip
	for _, ip := range loadBalancer.Ips {
		ips = append(ips, adaptIp(ip))
	}

	return &model.LoadBalancer{
		Id:        basetypes.NewStringValue(loadBalancer.Id),
		Type:      adaptInstanceType(loadBalancer.Type),
		Resources: adaptResources(loadBalancer.Resources),
		Region:    *adaptRegion(loadBalancer.Region),
		Reference: shared.AdaptNullableStringToStringValue(loadBalancer.Reference),
		State:     basetypes.NewStringValue(string(loadBalancer.State)),
		Contract:  adaptContract(loadBalancer.Contract),
		StartedAt: shared.AdaptNullableTimeToStringValue(loadBalancer.StartedAt),
		Ips:       ips,
		LoadBalancerConfiguration: shared.AdaptNullableDomainEntityToDatasourceModel(
			loadBalancer.Configuration,
			adaptLoadBalancerConfiguration,
		),
		PrivateNetwork: shared.AdaptNullableDomainEntityToDatasourceModel(
			loadBalancer.PrivateNetwork,
			adaptPrivateNetwork,
		),
	}
}

func adaptLoadBalancerConfiguration(configuration public_cloud.LoadBalancerConfiguration) *model.LoadBalancerConfiguration {
	return &model.LoadBalancerConfiguration{
		Balance: basetypes.NewStringValue(configuration.Balance.String()),
		StickySession: shared.AdaptNullableDomainEntityToDatasourceModel(
			configuration.StickySession,
			adaptStickySession,
		),
		XForwardedFor: basetypes.NewBoolValue(configuration.XForwardedFor),
		IdleTimeout:   basetypes.NewInt64Value(int64(configuration.IdleTimeout)),
	}
}

func adaptStickySession(stickySession public_cloud.StickySession) *model.StickySession {
	return &model.StickySession{
		Enabled:     basetypes.NewBoolValue(stickySession.Enabled),
		MaxLifeTime: basetypes.NewInt64Value(int64(stickySession.MaxLifeTime)),
	}
}

func adaptPrivateNetwork(privateNetwork public_cloud.PrivateNetwork) *model.PrivateNetwork {
	return &model.PrivateNetwork{
		Id:     basetypes.NewStringValue(privateNetwork.Id),
		Status: basetypes.NewStringValue(privateNetwork.Status),
		Subnet: basetypes.NewStringValue(privateNetwork.Subnet),
	}
}

func adaptIso(iso public_cloud.Iso) *model.Iso {
	return &model.Iso{
		Id:   basetypes.NewStringValue(iso.Id),
		Name: basetypes.NewStringValue(iso.Name),
	}
}

func adaptIp(ip public_cloud.Ip) model.Ip {
	return model.Ip{
		Ip:            basetypes.NewStringValue(ip.Ip),
		PrefixLength:  basetypes.NewStringValue(ip.PrefixLength),
		Version:       basetypes.NewInt64Value(int64(ip.Version)),
		NullRouted:    basetypes.NewBoolValue(ip.NullRouted),
		MainIp:        basetypes.NewBoolValue(ip.MainIp),
		NetworkType:   basetypes.NewStringValue(string(ip.NetworkType)),
		ReverseLookup: shared.AdaptNullableStringToStringValue(ip.ReverseLookup),
		Ddos: shared.AdaptNullableDomainEntityToDatasourceModel(
			ip.Ddos,
			adaptDdos,
		),
	}
}

func adaptDdos(ddos public_cloud.Ddos) *model.Ddos {
	return &model.Ddos{
		DetectionProfile: basetypes.NewStringValue(ddos.DetectionProfile),
		ProtectionType:   basetypes.NewStringValue(ddos.ProtectionType),
	}
}

func adaptStorageSize(storageSize public_cloud.StorageSize) *model.StorageSize {
	return &model.StorageSize{
		Size: basetypes.NewFloat64Value(storageSize.Size),
		Unit: basetypes.NewStringValue(storageSize.Unit),
	}
}

func adaptInstanceType(domainInstanceType public_cloud.InstanceType) model.InstanceType {
	instanceType := model.InstanceType{
		Name:      basetypes.NewStringValue(domainInstanceType.Name),
		Resources: adaptResources(domainInstanceType.Resources),
		Prices:    adaptPrices(domainInstanceType.Prices),
	}

	if domainInstanceType.StorageTypes != nil {
		instanceType.StorageTypes = domainInstanceType.StorageTypes.ToArray()
	}

	return instanceType
}

func adaptPrices(prices public_cloud.Prices) model.Prices {
	return model.Prices{
		Currency:       basetypes.NewStringValue(prices.Currency),
		CurrencySymbol: basetypes.NewStringValue(prices.CurrencySymbol),
		Compute:        adaptPrice(prices.Compute),
		Storage:        adaptStorage(prices.Storage),
	}
}

func adaptPrice(price public_cloud.Price) model.Price {
	return model.Price{
		HourlyPrice:  basetypes.NewStringValue(price.HourlyPrice),
		MonthlyPrice: basetypes.NewStringValue(price.MonthlyPrice),
	}
}

func adaptStorage(storage public_cloud.Storage) model.Storage {
	return model.Storage{
		Local:   adaptPrice(storage.Local),
		Central: adaptPrice(storage.Central),
	}
}

func adaptRegion(region public_cloud.Region) *model.Region {
	return &model.Region{
		Name:     basetypes.NewStringValue(region.Name),
		Location: basetypes.NewStringValue(region.Location),
	}
}
