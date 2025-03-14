---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ccx_datastore Resource - terraform-provider-ccx"
subcategory: ""
description: |-
  
---

# ccx_datastore (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_provider` (String) Cloud provider name
- `cloud_region` (String) The region to set up the datastore
- `db_vendor` (String) Database Vendor
- `instance_size` (String) Instance type/flavor to use
- `name` (String) The name of the datastore

### Optional

- `db_version` (String) Database Version
- `firewall` (Block List) FirewallRule rules to allow (see [below for nested schema](#nestedblock--firewall))
- `maintenance_day_of_week` (Number) Day of the week to run the maintenance. 1-7, 1 is Monday
- `maintenance_end_hour` (Number) Hour of the day to end the maintenance. 0-23. Must be start_hour + 2
- `maintenance_start_hour` (Number) Hour of the day to start the maintenance. 0-23
- `network_az` (List of String) Network availability zones
- `network_ha_enabled` (Boolean) High availability enabled or not
- `network_vpc_uuid` (String) VPC to use
- `notifications_emails` (List of String) List of email addresses to send notifications to
- `notifications_enabled` (Boolean) Enable or disable notifications. Default is false
- `parameter_group` (String) Parameter group ID to use
- `size` (Number) The size of the datastore ( int64 ). 1 or 3 nodes.
- `tags` (List of String) An optional list of tags
- `type` (String) Replication type of the datastore
- `volume_iops` (Number) Volume IOPS
- `volume_size` (Number) Volume size
- `volume_type` (String) Volume type

### Read-Only

- `dbname` (String) Database name
- `id` (String) The ID of this resource.
- `password` (String) Password to connect to the datastore
- `primary_dsn` (String) DSN to the primary host(s)
- `primary_url` (String) URL to the primary host(s)
- `replica_dsn` (String) DSN to the replica host(s)
- `replica_url` (String) URL to the replica host(s)
- `username` (String) Username to connect to the datastore

<a id="nestedblock--firewall"></a>
### Nested Schema for `firewall`

Required:

- `source` (String) CIDR source for the firewall rule

Optional:

- `description` (String) Description of this firewall rule

Read-Only:

- `id` (String)
