# Code architecture

## Packages

Each package corresponds to a product group, i.e.: `publiccloud`

## Data Sources

### Data Source files

All code pertaining to Data Sources, including models, belongs in the dataSource
file.
The format for this file is `<ENDPOINT>_data_source.go`.
For instances the filename would be `instances_data_source.go`

### Data Source structs

- Data Source structs should adhere to the following naming convention:
`<NAME>DataSource`. For `instances` this would be `instancesDataSource`
- Implement `utils.DataSourceAPI` in the struct to support `Metadata()`
  & `Configuration()` functions.

### Data Source models

Data Source model structs should adhere to the following convention:
`<MODEL_NAME>DataSourceModel`.
The `instances` data model would be named `instancesDataSourceModel`.

#### Data Source adaptation functions

Adapt functions to convert SDK models to Data Source models have the following
naming convention: `adapt<SDK_MODEL_NAME>To<DataSource_MODEL_NAME>DataSource`.
The function to adapt an SDK Instance to an Instance Data Source would be
named `adaptInstanceToInstanceData Source`.

## Resources

### Resource files

All code pertaining to resources, including models, belongs in the resource
file.
The format for this file is `<ENDPOINT>_resource.go`.
For instances the filename would be `instance_resource.go`

### Resource structs

- Resource structs should adhere to the following naming convention:
`<NAME>Resource`. For `instance` this would be `instanceResource`.
- Implement `utils.ResourceAPI` in the struct to support `Metadata()`
& `Configuration()` functions.

### Resource models

Resource model structs should adhere to the following convention:
`<MODEL_NAME>ResourceModel`.
The `instance` data model would be named `instanceResourceModel`

#### Resource adaptation functions

Adapt functions to convert SDK models to resource models have the following
naming convention: `adapt<SDK_MODEL_NAME>To<RESOURCE_MODEL_NAME>Resource`.
The function to adapt an SDK Instance to an Instance Resource would be
named `adaptInstanceToInstanceResource`.

## Validators

As validators are often shared between resources, they belong in the `validators.go`
file.

## SDK

Where possible, use the SDK getters.
Instead of `sdkInstance.Id` use `sdkInstance.GetId()`.

## Metadata

For both resources & data sources,
the following naming conventions apply for type names:

- The type names must have a `leaseweb` prefix
- `response.TypeName` names must be lowercase
- an underscore must replace non-alphabetic characters

`Public Cloud Load Balancers` thus becomes `leaseweb_public_cloud_load_balancers`

### Backwards compatibility

To maintain backwards compatibility,
the following resources & data sources do not adhere to the rules:

- `leseweb_dedicated_server`
- `leseweb_dedicated_servers`
- `leseweb_dedicated_server_notification_setting_bandwidth`
- `leseweb_dedicated_server_notification_setting_datatraffic`

## Tests

Acceptance tests should be placed in [internal/provider/provider_test.go](internal/provider/provider_test.go).
The test functions themselves must adhere to the
`TestAcc<PRODUCT_NAME><END_POINT><TYPE>` format. An acceptance test for a
Public Cloud instance resource should be called `TestAccPublicCloudInstanceResource` and a test for a
Public Cloud instance data source should be called `TestAccPublicCloudInstanceDataSource`.

## Variables

Where possible use `plan`,`state` & `config` as names for variables that reference terraform plan, state & config.
