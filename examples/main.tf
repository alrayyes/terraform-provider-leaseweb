terraform {
  required_providers {
    leaseweb = {
      version = "0.0.1"
      source  = "git.ocom.com/infra/leaseweb"
    }
  }
}

provider "leaseweb" {}

resource "leaseweb_dedicatedserver" "my-test" {
  # reference = "web01"
  # reverse_lookup = "web02.example.com"
  # dhcp_lease = "https://boot.netboot.xyz"
  # powered_on = true
  # main_ip_nulled = false
}

resource "leaseweb_dedicatedserver_notification_setting_bandwidth" "alert" {
  dedicated_server_id = leaseweb_dedicatedserver.my-test.id
  frequency = "DAILY"
  threshold = 1.5
  unit = "Gbps"
}

resource "leaseweb_dedicatedserver_notification_setting_datatraffic" "alert" {
  dedicated_server_id = leaseweb_dedicatedserver.my-test.id
  frequency = "WEEKLY"
  threshold = 2
  unit = "TB"
}

data "leaseweb_operating_systems" "all_os" {
}

locals {
  latest_ubuntu_os_id = reverse(sort([
    for id in data.leaseweb_operating_systems.all_os.ids: id
    if length(regexall("^UBUNTU_.*", id)) > 0
  ]))[0]
}

output "latest_ubuntu_os_name" {
  value = data.leaseweb_operating_systems.all_os.names[local.latest_ubuntu_os_id]
}
