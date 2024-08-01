package public_cloud

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/services/errors"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

// Service fulfills the contract for ports.PublicCloudService.
type Service struct {
	publicCloudRepository ports.PublicCloudRepository
}

func (srv Service) GetAllInstances(ctx context.Context) (
	domain.Instances,
	*errors.ServiceError,
) {
	var detailedInstances domain.Instances
	resultChan := make(chan domain.Instance)
	errorChan := make(chan *errors.ServiceError)

	instances, err := srv.publicCloudRepository.GetAllInstances(ctx)
	if err != nil {
		return domain.Instances{}, errors.NewFromRepositoryError(
			"GetAllInstances",
			*err,
		)
	}

	for _, instance := range instances {
		go func(id value_object.Uuid) {
			detailedInstance, err := srv.GetInstance(id, ctx)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- *detailedInstance
		}(instance.Id)
	}

	for i := 0; i < len(instances); i++ {
		select {
		case err := <-errorChan:
			return domain.Instances{}, err
		case res := <-resultChan:
			detailedInstances = append(detailedInstances, res)
		}
	}

	return detailedInstances, nil
}

func (srv Service) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.Instance, *errors.ServiceError) {
	instance, err := srv.publicCloudRepository.GetInstance(id, ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("GetInstance", *err)
	}

	return srv.populateMissingInstanceAttributes(*instance, ctx)
}

func (srv Service) CreateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *errors.ServiceError) {
	createdInstance, err := srv.publicCloudRepository.CreateInstance(
		instance,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError("CreateInstance", *err)
	}

	// call GetInstance as createdInstance is created from instance and not instanceDetails
	return srv.GetInstance(createdInstance.Id, ctx)
}

func (srv Service) UpdateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *errors.ServiceError) {
	updatedInstance, err := srv.publicCloudRepository.UpdateInstance(
		instance,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError("UpdateInstance", *err)
	}

	return srv.populateMissingInstanceAttributes(*updatedInstance, ctx)
}

func (srv Service) DeleteInstance(
	id value_object.Uuid,
	ctx context.Context,
) *errors.ServiceError {
	err := srv.publicCloudRepository.DeleteInstance(id, ctx)
	if err != nil {
		return errors.NewError("DeleteInstance", err)
	}

	return nil
}

func (srv Service) GetAvailableInstanceTypesForUpdate(
	id value_object.Uuid,
	ctx context.Context,
) (domain.InstanceTypes, *errors.ServiceError) {
	instanceTypes, err := srv.publicCloudRepository.GetAvailableInstanceTypesForUpdate(
		id,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAvailableInstanceTypesForUpdate",
			*err,
		)
	}

	return instanceTypes, nil
}

func (srv Service) GetRegions(ctx context.Context) (
	domain.Regions,
	*errors.ServiceError,
) {
	regions, err := srv.publicCloudRepository.GetRegions(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("GetRegions", *err)
	}

	return regions, nil
}

// Populate instance with autoScalingGroupDetails & loadBalancerDetails.
func (srv Service) populateMissingInstanceAttributes(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *errors.ServiceError) {
	// Get autoScalingGroupDetails.
	if instance.AutoScalingGroup != nil {
		autoScalingGroup, err := srv.publicCloudRepository.GetAutoScalingGroup(
			instance.AutoScalingGroup.Id,
			ctx,
		)
		if err != nil {
			return nil, errors.NewFromRepositoryError(
				"populateMissingInstanceAttributes",
				*err,
			)
		}

		// Get loadBalancerDetails.
		if autoScalingGroup.LoadBalancer != nil {
			loadBalancer, err := srv.publicCloudRepository.GetLoadBalancer(
				autoScalingGroup.LoadBalancer.Id,
				ctx,
			)
			if err != nil {
				return nil, errors.NewFromRepositoryError(
					"populateMissingInstanceAttributes",
					*err,
				)
			}
			autoScalingGroup.LoadBalancer = loadBalancer
		}

		instance.AutoScalingGroup = autoScalingGroup
	}

	return &instance, nil
}

func (srv Service) GetAvailableInstanceTypesForRegion(
	region string,
	ctx context.Context,
) (domain.InstanceTypes, *errors.ServiceError) {
	instanceTypes, err := srv.publicCloudRepository.GetInstanceTypesForRegion(
		region,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"populateMissingInstanceAttributes",
			*err,
		)
	}

	return instanceTypes, nil
}

func New(publicCloudRepository ports.PublicCloudRepository) Service {
	return Service{publicCloudRepository: publicCloudRepository}
}