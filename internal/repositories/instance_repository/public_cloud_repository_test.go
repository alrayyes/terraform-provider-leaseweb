package instance_repository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

var (
	_ publicCloudApi = &publicCloudApiSpy{}
)

var sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

type publicCloudApiSpy struct {
	instance                        *publicCloud.InstanceDetails
	instanceList                    []publicCloud.Instance
	autoScalingGroup                *publicCloud.AutoScalingGroupDetails
	loadBalancer                    *publicCloud.LoadBalancerDetails
	launchedInstance                *publicCloud.Instance
	getInstanceListExecuteError     error
	getInstanceExecuteError         error
	getAutoScalingGroupExecuteError error
	getLoadBalancerExecuteError     error
	launchInstanceExecuteError      error
	updateInstanceExecuteError      error
}

func (a publicCloudApiSpy) LaunchInstance(ctx context.Context) publicCloud.ApiLaunchInstanceRequest {
	return publicCloud.ApiLaunchInstanceRequest{}
}

func (a publicCloudApiSpy) LaunchInstanceExecute(r publicCloud.ApiLaunchInstanceRequest) (
	*publicCloud.Instance,
	*http.Response,
	error,
) {
	return a.launchedInstance, nil, a.launchInstanceExecuteError
}

func (a publicCloudApiSpy) UpdateInstance(
	ctx context.Context,
	instanceId string,
) publicCloud.ApiUpdateInstanceRequest {
	return publicCloud.ApiUpdateInstanceRequest{}
}

func (a publicCloudApiSpy) UpdateInstanceExecute(r publicCloud.ApiUpdateInstanceRequest) (
	*publicCloud.InstanceDetails,
	*http.Response,
	error,
) {
	return a.instance, nil, a.updateInstanceExecuteError
}

func (a publicCloudApiSpy) GetInstanceList(ctx context.Context) publicCloud.ApiGetInstanceListRequest {
	return publicCloud.ApiGetInstanceListRequest{}
}

func (a publicCloudApiSpy) GetInstanceListExecute(r publicCloud.ApiGetInstanceListRequest) (
	*publicCloud.GetInstanceListResult,
	*http.Response,
	error,
) {
	return &publicCloud.GetInstanceListResult{Instances: a.instanceList}, nil, a.getInstanceListExecuteError
}

func (a publicCloudApiSpy) GetAutoScalingGroup(
	ctx context.Context,
	autoScalingGroupId string,
) publicCloud.ApiGetAutoScalingGroupRequest {
	return publicCloud.ApiGetAutoScalingGroupRequest{}

}

func (a publicCloudApiSpy) GetAutoScalingGroupExecute(r publicCloud.ApiGetAutoScalingGroupRequest) (
	*publicCloud.AutoScalingGroupDetails,
	*http.Response,
	error,
) {

	return a.autoScalingGroup, nil, a.getAutoScalingGroupExecuteError
}

func (a publicCloudApiSpy) GetLoadBalancer(
	ctx context.Context,
	loadBalancerId string,
) publicCloud.ApiGetLoadBalancerRequest {
	return publicCloud.ApiGetLoadBalancerRequest{}
}

func (a publicCloudApiSpy) GetLoadBalancerExecute(r publicCloud.ApiGetLoadBalancerRequest) (
	*publicCloud.LoadBalancerDetails,
	*http.Response,
	error,
) {
	return a.loadBalancer, nil, a.getLoadBalancerExecuteError
}

func (a publicCloudApiSpy) GetInstance(
	ctx context.Context,
	instanceId string,
) publicCloud.ApiGetInstanceRequest {
	return publicCloud.ApiGetInstanceRequest{}
}

func (a publicCloudApiSpy) GetInstanceExecute(r publicCloud.ApiGetInstanceRequest) (
	*publicCloud.InstanceDetails,
	*http.Response,
	error,
) {
	return a.instance, nil, a.getInstanceExecuteError
}

func TestNewPublicCloudRepository(t *testing.T) {
	t.Run("token is set properly", func(t *testing.T) {
		got := NewPublicCloudRepository("token", Optional{})

		assert.Equal(t, "token", got.token)
	})
}

func TestPublicCloudRepository_authContext(t *testing.T) {
	publicCloudRepository := NewPublicCloudRepository("token", Optional{})
	got := publicCloudRepository.authContext(context.TODO()).Value(publicCloud.ContextAPIKeys)

	assert.Equal(
		t,
		map[string]publicCloud.APIKey{"X-LSW-Auth": {Key: "token", Prefix: ""}},
		got,
	)
}

func TestPublicCloudRepository_GetInstance(t *testing.T) {
	t.Run("expected instance entity is returned", func(t *testing.T) {
		id := uuid.New()
		convertedInstanceId, _ := uuid.NewUUID()

		apiSpy := publicCloudApiSpy{instance: &publicCloud.InstanceDetails{Id: id.String()}}

		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertInstance: func(
				skInstance publicCloud.InstanceDetails,
				autoScalingGroup *entity.AutoScalingGroup,
			) (*entity.Instance, error) {
				assert.Equal(
					t,
					id.String(),
					skInstance.GetId(),
					"sdkInstance is converted",
				)

				return &entity.Instance{Id: convertedInstanceId}, nil
			},
		}

		got, err := publicCloudRepository.GetInstance(id, context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, convertedInstanceId, got.Id)
	})

	t.Run(
		"error is returned if instance cannot be retrieved from the sdk",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{getInstanceExecuteError: errors.New("error getting instance")}

			PublicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := PublicCloudRepository.GetInstance(uuid.New(), context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting instance")
		},
	)

	t.Run("expected autoScalingGroup is set", func(t *testing.T) {
		autoScalingGroupId := uuid.New()
		convertedAutoScalingGroupId := uuid.New()

		apiSpy := publicCloudApiSpy{
			instance: &publicCloud.InstanceDetails{
				AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{
					Id: autoScalingGroupId.String()},
				),
			},
			autoScalingGroup: &publicCloud.AutoScalingGroupDetails{Id: autoScalingGroupId.String()},
		}

		PublicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertInstance: func(
				skInstance publicCloud.InstanceDetails,
				autoScalingGroup *entity.AutoScalingGroup,
			) (*entity.Instance, error) {
				assert.Equal(
					t,
					convertedAutoScalingGroupId,
					autoScalingGroup.Id,
					"autoScalingGroup is passed on to convertInstance",
				)

				return &entity.Instance{AutoScalingGroup: &entity.AutoScalingGroup{Id: convertedAutoScalingGroupId}}, nil
			},
			convertAutoScalingGroup: func(
				sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
				loadBalancer *entity.LoadBalancer,
			) (*entity.AutoScalingGroup, error) {
				assert.Equal(
					t,
					autoScalingGroupId.String(),
					sdkAutoScalingGroup.GetId(),
					"sdkAutoScalingGroup is converted",
				)
				return &entity.AutoScalingGroup{Id: convertedAutoScalingGroupId}, nil
			},
		}

		got, err := PublicCloudRepository.GetInstance(uuid.New(), context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, convertedAutoScalingGroupId, got.AutoScalingGroup.Id)
	})

	t.Run("error is returned if autoScalingGroup uuid is invalid", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			instance: &publicCloud.InstanceDetails{
				AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{
					Id: "tralala"},
				),
			},
		}

		PublicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

		_, err := PublicCloudRepository.GetInstance(uuid.New(), context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot convert string to uuid")
	})

	t.Run("error is returned if autoScalingGroup cannot be retrieved", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			instance: &publicCloud.InstanceDetails{
				AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{
					Id: uuid.New().String()},
				),
			},
			getAutoScalingGroupExecuteError: errors.New("error getting autoScalingGroup"),
		}

		PublicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

		_, err := PublicCloudRepository.GetInstance(uuid.New(), context.TODO())

		assert.Error(t, err)
		assert.ErrorContains(t, err, "error getting autoScalingGroup")
	})
}

func TestPublicCLoudRepository_GetLoadBalancer(t *testing.T) {
	t.Run("expected loadBalancer entity is returned", func(t *testing.T) {
		id := uuid.New()
		convertedId := uuid.New()

		apiSpy := publicCloudApiSpy{
			loadBalancer: &publicCloud.LoadBalancerDetails{Id: id.String()},
		}

		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertLoadBalancer: func(
				sdkLoadBalancer publicCloud.LoadBalancerDetails,
			) (*entity.LoadBalancer, error) {
				assert.Equal(
					t,
					id.String(),
					sdkLoadBalancer.Id,
					"sdkLoadBalancer is passed on to convertLoadBalancer",
				)

				return &entity.LoadBalancer{Id: convertedId}, nil
			},
		}

		got, err := publicCloudRepository.GetLoadBalancer(id, context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, convertedId, got.Id)
	})

	t.Run(
		"error is returned when loadBalancer cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				getLoadBalancerExecuteError: errors.New("error getting loadBalancer"),
			}

			PublicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := PublicCloudRepository.GetLoadBalancer(uuid.New(), context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting loadBalancer")
		},
	)

	t.Run(
		"error is returned if loadBalancer cannot be converted",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				loadBalancer: &publicCloud.LoadBalancerDetails{Id: uuid.New().String()},
			}

			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertLoadBalancer: func(
					sdkLoadBalancer publicCloud.LoadBalancerDetails,
				) (*entity.LoadBalancer, error) {

					return nil, errors.New("conversion error")
				},
			}

			_, err := publicCloudRepository.GetLoadBalancer(
				uuid.New(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "conversion error")
		},
	)
}

func TestPublicCloudRepository_GetAutoScalingGroup(t *testing.T) {
	t.Run(
		"expected autoScalingGroup entity is returned",
		func(t *testing.T) {
			id := uuid.New()
			convertedId := uuid.New()
			apiSpy := publicCloudApiSpy{
				autoScalingGroup: &publicCloud.AutoScalingGroupDetails{Id: id.String()},
			}

			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertAutoScalingGroup: func(
					sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
					loadBalancer *entity.LoadBalancer,
				) (*entity.AutoScalingGroup, error) {
					assert.Equal(
						t,
						id.String(),
						sdkAutoScalingGroup.Id,
						"sdkLoadBalancer is passed on to convertLoadBalancer",
					)

					return &entity.AutoScalingGroup{Id: convertedId}, nil
				},
			}

			got, err := publicCloudRepository.GetAutoScalingGroup(id, context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, convertedId, got.Id)
		},
	)

	t.Run(
		"return error if autoScalingGroup cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				getAutoScalingGroupExecuteError: errors.New("error getting autoScalingGroup"),
			}

			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAutoScalingGroup(
				uuid.New(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting autoScalingGroup")
		},
	)

	t.Run("return error if loadBalancer id is invalid", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			autoScalingGroup: &publicCloud.AutoScalingGroupDetails{
				LoadBalancer: *publicCloud.NewNullableLoadBalancer(&publicCloud.LoadBalancer{Id: "tralala"}),
			},
		}

		publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

		_, err := publicCloudRepository.GetAutoScalingGroup(
			uuid.New(),
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot convert string to uuid")
	},
	)

	t.Run(
		"return error if loadBalancer cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				autoScalingGroup: &publicCloud.AutoScalingGroupDetails{
					LoadBalancer: *publicCloud.NewNullableLoadBalancer(&publicCloud.LoadBalancer{
						Id: uuid.New().String()},
					),
				},
				getLoadBalancerExecuteError: errors.New("error getting loadBalancer"),
			}

			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAutoScalingGroup(
				uuid.New(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting loadBalancer")
		},
	)

	t.Run(
		"return error if autoScalingGroup cannot be converted",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				autoScalingGroup: &publicCloud.AutoScalingGroupDetails{},
			}

			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertAutoScalingGroup: func(
					sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
					loadBalancer *entity.LoadBalancer,
				) (*entity.AutoScalingGroup, error) {
					return nil, errors.New("conversion error")
				},
			}

			_, err := publicCloudRepository.GetAutoScalingGroup(
				uuid.New(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "conversion error")
		},
	)

}

func TestPublicCloudRepository_GetAllInstances(t *testing.T) {
	t.Run("expected instances entity is returned", func(t *testing.T) {
		id := uuid.New()

		apiSpy := publicCloudApiSpy{
			instanceList: []publicCloud.Instance{{Id: id.String()}},
			instance:     &publicCloud.InstanceDetails{Id: id.String()},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertInstance: func(
				sdkInstance publicCloud.InstanceDetails,
				sdkAutoScalingGroup *entity.AutoScalingGroup,
			) (*entity.Instance, error) {
				return &entity.Instance{Id: id}, nil
			},
		}

		got, err := publicCloudRepository.GetAllInstances(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, id, got[0].Id)
	})

	t.Run(
		"return error when instances cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				getInstanceListExecuteError: errors.New("error getting instances"),
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting instances")
		},
	)

	t.Run(
		"return error when instance id cannot be parsed",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				instanceList: []publicCloud.Instance{{Id: "tralala"}},
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"return error when instance cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				instanceList:            []publicCloud.Instance{{Id: uuid.New().String()}},
				getInstanceExecuteError: errors.New("error getting instance"),
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting instance")
		},
	)
}

func TestPublicCloudRepository_CreateInstance(t *testing.T) {
	t.Run("expected instance entity is created", func(t *testing.T) {
		id := uuid.New()
		convertedId := uuid.New()

		apiSpy := publicCloudApiSpy{
			launchedInstance: &publicCloud.Instance{Id: id.String()},
			instance:         &publicCloud.InstanceDetails{Id: id.String()},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertInstance: func(sdkInstance publicCloud.InstanceDetails, sdkAutoScalingGroup *entity.AutoScalingGroup) (*entity.Instance, error) {
				return &entity.Instance{Id: convertedId}, nil
			},
		}

		marketAppId := "marketAppId"
		reference := "reference"
		sshKeyValueObject, _ := value_object.NewSshKey(sshKey)

		instance := entity.NewCreateInstance(
			"region",
			"lsw.m3.large",
			enum.RootDiskStorageTypeCentral,
			enum.Almalinux864Bit,
			enum.ContractTypeMonthly,
			enum.ContractTermSix,
			enum.ContractBillingFrequencyThree,
			entity.OptionalCreateInstanceValues{
				MarketAppId: &marketAppId,
				Reference:   &reference,
				SshKey:      sshKeyValueObject,
			},
		)

		got, err := publicCloudRepository.CreateInstance(
			instance,
			context.TODO(),
		)

		assert.NoError(t, err)
		assert.Equal(t, convertedId, got.Id)
	})

	t.Run("invalid instanceType returns error", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			launchedInstance: &publicCloud.Instance{Id: uuid.New().String()},
			instance:         &publicCloud.InstanceDetails{Id: uuid.New().String()},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
		}

		_, err := publicCloudRepository.CreateInstance(
			entity.Instance{Type: "tralala"},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid rootDiskStorageType returns error", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			launchedInstance: &publicCloud.Instance{Id: uuid.New().String()},
			instance:         &publicCloud.InstanceDetails{Id: uuid.New().String()},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
		}

		_, err := publicCloudRepository.CreateInstance(
			entity.Instance{Type: "lsw.m3.large", RootDiskStorageType: "tralala"},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid imageId returns error", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			launchedInstance: &publicCloud.Instance{Id: uuid.New().String()},
			instance:         &publicCloud.InstanceDetails{Id: uuid.New().String()},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
		}

		_, err := publicCloudRepository.CreateInstance(
			entity.Instance{
				Type:                "lsw.m3.large",
				RootDiskStorageType: enum.RootDiskStorageTypeCentral,
				Image:               entity.Image{Id: "tralala"},
			},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractType returns error", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			launchedInstance: &publicCloud.Instance{Id: uuid.New().String()},
			instance:         &publicCloud.InstanceDetails{Id: uuid.New().String()},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
		}

		_, err := publicCloudRepository.CreateInstance(
			entity.Instance{
				Type:                "lsw.m3.large",
				RootDiskStorageType: enum.RootDiskStorageTypeCentral,
				Image:               entity.Image{Id: enum.Ubuntu200464Bit},
				Contract:            entity.Contract{Type: "tralala"},
			},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid contractTerm returns error", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			launchedInstance: &publicCloud.Instance{Id: uuid.New().String()},
			instance:         &publicCloud.InstanceDetails{Id: uuid.New().String()},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
		}

		_, err := publicCloudRepository.CreateInstance(
			entity.Instance{
				Type:                "lsw.m3.large",
				RootDiskStorageType: enum.RootDiskStorageTypeCentral,
				Image:               entity.Image{Id: enum.Ubuntu200464Bit},
				Contract: entity.Contract{
					Type: enum.ContractTypeMonthly,
					Term: 55,
				},
			},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run("invalid billingFrequency returns error", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			launchedInstance: &publicCloud.Instance{Id: uuid.New().String()},
			instance:         &publicCloud.InstanceDetails{Id: uuid.New().String()},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
		}

		_, err := publicCloudRepository.CreateInstance(
			entity.Instance{
				Type:                "lsw.m3.large",
				RootDiskStorageType: enum.RootDiskStorageTypeCentral,
				Image:               entity.Image{Id: enum.Ubuntu200464Bit},
				Contract: entity.Contract{
					Type:             enum.ContractTypeMonthly,
					Term:             enum.ContractTermThree,
					BillingFrequency: 55,
				},
			},
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "55")
	})

	t.Run(
		"error is returned when instance cannot be launched in sdk",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				launchInstanceExecuteError: errors.New("some error"),
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			instance := entity.NewCreateInstance(
				"region",
				"lsw.m3.large",
				enum.RootDiskStorageTypeCentral,
				enum.Almalinux864Bit,
				enum.ContractTypeMonthly,
				enum.ContractTermSix,
				enum.ContractBillingFrequencyThree,
				entity.OptionalCreateInstanceValues{},
			)

			_, err := publicCloudRepository.CreateInstance(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned when id of launched instance is invalid",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				launchedInstance: &publicCloud.Instance{Id: "tralala"},
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			instance := entity.NewCreateInstance(
				"region",
				"lsw.m3.large",
				enum.RootDiskStorageTypeCentral,
				enum.Almalinux864Bit,
				enum.ContractTypeMonthly,
				enum.ContractTermSix,
				enum.ContractBillingFrequencyThree,
				entity.OptionalCreateInstanceValues{},
			)

			_, err := publicCloudRepository.CreateInstance(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"error is returned when instanceDetails cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				launchedInstance:        &publicCloud.Instance{Id: uuid.New().String()},
				getInstanceExecuteError: errors.New("some error"),
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			instance := entity.NewCreateInstance(
				"region",
				"lsw.m3.large",
				enum.RootDiskStorageTypeCentral,
				enum.Almalinux864Bit,
				enum.ContractTypeMonthly,
				enum.ContractTermSix,
				enum.ContractBillingFrequencyThree,
				entity.OptionalCreateInstanceValues{},
			)

			_, err := publicCloudRepository.CreateInstance(
				instance,
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestPublicCloudRepository_UpdateInstance(t *testing.T) {
	t.Run("expected instance entity is updated", func(t *testing.T) {
		publicCloudRepository := PublicCloudRepository{}
		got, err := publicCloudRepository.UpdateInstance(
			entity.Instance{},
			context.TODO(),
		)

		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func TestPublicCloudRepository_DeleteInstance(t *testing.T) {
	t.Run("expected instance entity is deleted", func(t *testing.T) {
		publicCloudRepository := PublicCloudRepository{}
		err := publicCloudRepository.DeleteInstance(uuid.New(), context.TODO())

		assert.NoError(t, err)
	})
}
