---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "leaseweb_dedicated_server_control_panels Data Source - terraform-provider-leaseweb"
subcategory: ""
description: |-
  The dedicated_server_control_panels data source allows access to the list of
  control panels available for installation on a dedicated server.
---

# leaseweb_dedicated_server_control_panels (Data Source)

The `dedicated_server_control_panels` data source allows access to the list of
control panels available for installation on a dedicated server.

## Example Usage

```terraform
# Access all control panels available with Ubuntu 22.04
data "leaseweb_dedicated_server_control_panels" "ubuntu_cps" {
  operating_system_id = "UBUNTU_22_04_64BIT"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `operating_system_id` (String) Filter the list of control panels to return only the ones available to an operating system.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of String) List of the control panel IDs.
- `names` (Map of String) List of the control panel names.

