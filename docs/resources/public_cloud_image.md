---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "leaseweb_public_cloud_image Resource - leaseweb"
subcategory: ""
description: |-
  Once created, an image resource cannot be deleted via Terraform
---

# leaseweb_public_cloud_image (Resource)

Once created, an image resource cannot be deleted via Terraform

## Example Usage

```terraform
# Manage example Public Cloud image
resource "leaseweb_public_cloud_image" "example" {
  id   = "396a3299-1795-464b-aa10-e1f179db1926"
  name = "Custom image"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The id of the instance which the custom image is based on. The following rules apply:
  - instance exists for instanceId
  - instance has state *STOPPED*
  - instance has a maximum rootDiskSize of 100 GB
  - instance OS must not be *windows*
- `name` (String) Custom image name

### Read-Only

- `custom` (Boolean) Standard or Custom image
- `flavour` (String)
- `market_apps` (List of String)
- `region` (String)
- `state` (String)
- `storage_types` (List of String) The supported storage types for the instance type

## Import

Import is supported using the following syntax:

```shell
# Public Cloud image can be imported by specifying the identifier.
terraform import leaseweb_public_cloud_image.example ace712e9-a166-47f1-9065-4af0f7e7fce1
```