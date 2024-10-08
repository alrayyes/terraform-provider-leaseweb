---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "leaseweb_public_cloud_instances Data Source - leaseweb"
subcategory: ""
description: |-
  
---

# leaseweb_public_cloud_instances (Data Source)



## Example Usage

```terraform
# List all Public Cloud instances
data "leaseweb_public_cloud_instances" "all" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `instances` (Attributes List) (see [below for nested schema](#nestedatt--instances))

<a id="nestedatt--instances"></a>
### Nested Schema for `instances`

Read-Only:

- `auto_scaling_group` (Attributes) (see [below for nested schema](#nestedatt--instances--auto_scaling_group))
- `contract` (Attributes) (see [below for nested schema](#nestedatt--instances--contract))
- `has_private_network` (Boolean)
- `has_public_ipv4` (Boolean)
- `has_user_data` (Boolean)
- `id` (String) The instance unique identifier
- `image` (Attributes) (see [below for nested schema](#nestedatt--instances--image))
- `ips` (Attributes List) (see [below for nested schema](#nestedatt--instances--ips))
- `iso` (Attributes) (see [below for nested schema](#nestedatt--instances--iso))
- `market_app_id` (String) Market App ID
- `private_network` (Attributes) (see [below for nested schema](#nestedatt--instances--private_network))
- `product_type` (String) The product type
- `reference` (String) The identifying name set to the instance
- `region` (Attributes) (see [below for nested schema](#nestedatt--instances--region))
- `resources` (Attributes) Available resources (see [below for nested schema](#nestedatt--instances--resources))
- `root_disk_size` (Number) The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances
- `root_disk_storage_type` (String) The root disk's storage type
- `started_at` (String) Date and time when the instance was started for the first time, right after launching it
- `state` (String) The instance's current state
- `type` (Attributes) (see [below for nested schema](#nestedatt--instances--type))

<a id="nestedatt--instances--auto_scaling_group"></a>
### Nested Schema for `instances.auto_scaling_group`

Read-Only:

- `cooldown_time` (Number) Only for "CPU_BASED" auto scaling group. Cool-down time in seconds for new instances
- `cpu_threshold` (Number) Only for "CPU_BASED" auto scaling group. The target average CPU utilization for scaling
- `created_at` (String) Date and time when the Auto Scaling Group was created
- `desired_amount` (Number) Number of instances that should be running
- `ends_at` (String) Only for "SCHEDULED" auto scaling group. Date and time (UTC) that the instances need to be terminated
- `id` (String) The Auto Scaling Group unique identifier
- `maximum_amount` (Number) Only for "CPU_BASED" auto scaling group. The maximum number of instances that can be running
- `minimum_amount` (Number) The minimum number of instances that should be running
- `reference` (String) The identifying name set to the auto scaling group
- `region` (Attributes) (see [below for nested schema](#nestedatt--instances--auto_scaling_group--region))
- `starts_at` (String) Only for "SCHEDULED" auto scaling group. Date and time (UTC) that the instances need to be launched
- `state` (String) The Auto Scaling Group's current state.
- `type` (String) Auto Scaling Group type
- `updated_at` (String) Date and time when the Auto Scaling Group was updated
- `warmup_time` (Number) Only for "CPU_BASED" auto scaling group. Warm-up time in seconds for new instances

<a id="nestedatt--instances--auto_scaling_group--region"></a>
### Nested Schema for `instances.auto_scaling_group.region`

Read-Only:

- `location` (String) The city where the region is located
- `name` (String)



<a id="nestedatt--instances--contract"></a>
### Nested Schema for `instances.contract`

Read-Only:

- `billing_frequency` (Number) The billing frequency (in months). Valid options are 
  - *0*
  - *1*
  - *3*
  - *6*
  - *12*
- `created_at` (String) Date when the contract was created
- `ends_at` (String)
- `renewals_at` (String) Date when the contract will be automatically renewed
- `state` (String)
- `term` (Number) Contract term (in months). Used only when type is *MONTHLY*. Valid options are 
  - *0*
  - *1*
  - *3*
  - *6*
  - *12*
- `type` (String) Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage


<a id="nestedatt--instances--image"></a>
### Nested Schema for `instances.image`

Read-Only:

- `architecture` (String)
- `created_at` (String)
- `custom` (Boolean)
- `family` (String)
- `flavour` (String)
- `id` (String) Image ID
- `market_apps` (List of String)
- `name` (String)
- `region` (Attributes) (see [below for nested schema](#nestedatt--instances--image--region))
- `state` (String)
- `state_reason` (String)
- `storage_size` (Attributes) (see [below for nested schema](#nestedatt--instances--image--storage_size))
- `storage_types` (List of String) The supported storage types
- `updated_at` (String)
- `version` (String)

<a id="nestedatt--instances--image--region"></a>
### Nested Schema for `instances.image.region`

Read-Only:

- `location` (String) The city where the region is located
- `name` (String)


<a id="nestedatt--instances--image--storage_size"></a>
### Nested Schema for `instances.image.storage_size`

Read-Only:

- `size` (Number) The storage size
- `unit` (String) The storage size unit



<a id="nestedatt--instances--ips"></a>
### Nested Schema for `instances.ips`

Read-Only:

- `ddos` (Attributes) (see [below for nested schema](#nestedatt--instances--ips--ddos))
- `ip` (String)
- `main_ip` (Boolean)
- `network_type` (String)
- `null_routed` (Boolean)
- `prefix_length` (String)
- `reverse_lookup` (String)
- `version` (Number)

<a id="nestedatt--instances--ips--ddos"></a>
### Nested Schema for `instances.ips.ddos`

Read-Only:

- `detection_profile` (String)
- `protection_type` (String)



<a id="nestedatt--instances--iso"></a>
### Nested Schema for `instances.iso`

Read-Only:

- `id` (String)
- `name` (String)


<a id="nestedatt--instances--private_network"></a>
### Nested Schema for `instances.private_network`

Read-Only:

- `id` (String)
- `status` (String)
- `subnet` (String)


<a id="nestedatt--instances--region"></a>
### Nested Schema for `instances.region`

Read-Only:

- `location` (String) The city where the region is located
- `name` (String)


<a id="nestedatt--instances--resources"></a>
### Nested Schema for `instances.resources`

Read-Only:

- `cpu` (Attributes) Number of cores (see [below for nested schema](#nestedatt--instances--resources--cpu))
- `memory` (Attributes) Total memory in GiB (see [below for nested schema](#nestedatt--instances--resources--memory))
- `private_network_speed` (Attributes) Private network speed in Gbps (see [below for nested schema](#nestedatt--instances--resources--private_network_speed))
- `public_network_speed` (Attributes) Public network speed in Gbps (see [below for nested schema](#nestedatt--instances--resources--public_network_speed))

<a id="nestedatt--instances--resources--cpu"></a>
### Nested Schema for `instances.resources.cpu`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--instances--resources--memory"></a>
### Nested Schema for `instances.resources.memory`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--instances--resources--private_network_speed"></a>
### Nested Schema for `instances.resources.private_network_speed`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--instances--resources--public_network_speed"></a>
### Nested Schema for `instances.resources.public_network_speed`

Read-Only:

- `unit` (String)
- `value` (Number)



<a id="nestedatt--instances--type"></a>
### Nested Schema for `instances.type`

Read-Only:

- `name` (String) Type name
- `prices` (Attributes) (see [below for nested schema](#nestedatt--instances--type--prices))
- `resources` (Attributes) Available resources (see [below for nested schema](#nestedatt--instances--type--resources))
- `storage_types` (List of String) The supported storage types

<a id="nestedatt--instances--type--prices"></a>
### Nested Schema for `instances.type.prices`

Read-Only:

- `compute` (Attributes) (see [below for nested schema](#nestedatt--instances--type--prices--compute))
- `currency` (String)
- `currency_symbol` (String)
- `storage` (Attributes) (see [below for nested schema](#nestedatt--instances--type--prices--storage))

<a id="nestedatt--instances--type--prices--compute"></a>
### Nested Schema for `instances.type.prices.compute`

Read-Only:

- `hourly_price` (String)
- `monthly_price` (String)


<a id="nestedatt--instances--type--prices--storage"></a>
### Nested Schema for `instances.type.prices.storage`

Read-Only:

- `central` (Attributes) (see [below for nested schema](#nestedatt--instances--type--prices--storage--central))
- `local` (Attributes) (see [below for nested schema](#nestedatt--instances--type--prices--storage--local))

<a id="nestedatt--instances--type--prices--storage--central"></a>
### Nested Schema for `instances.type.prices.storage.central`

Read-Only:

- `hourly_price` (String)
- `monthly_price` (String)


<a id="nestedatt--instances--type--prices--storage--local"></a>
### Nested Schema for `instances.type.prices.storage.local`

Read-Only:

- `hourly_price` (String)
- `monthly_price` (String)




<a id="nestedatt--instances--type--resources"></a>
### Nested Schema for `instances.type.resources`

Read-Only:

- `cpu` (Attributes) Number of cores (see [below for nested schema](#nestedatt--instances--type--resources--cpu))
- `memory` (Attributes) Total memory in GiB (see [below for nested schema](#nestedatt--instances--type--resources--memory))
- `private_network_speed` (Attributes) Private network speed in Gbps (see [below for nested schema](#nestedatt--instances--type--resources--private_network_speed))
- `public_network_speed` (Attributes) Public network speed in Gbps (see [below for nested schema](#nestedatt--instances--type--resources--public_network_speed))

<a id="nestedatt--instances--type--resources--cpu"></a>
### Nested Schema for `instances.type.resources.cpu`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--instances--type--resources--memory"></a>
### Nested Schema for `instances.type.resources.memory`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--instances--type--resources--private_network_speed"></a>
### Nested Schema for `instances.type.resources.private_network_speed`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--instances--type--resources--public_network_speed"></a>
### Nested Schema for `instances.type.resources.public_network_speed`

Read-Only:

- `unit` (String)
- `value` (Number)
