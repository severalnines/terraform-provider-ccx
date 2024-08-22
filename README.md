# terraform-provider-ccx

This is the Terraform Provider for the Severalnines CCX Database as a service platform.

- CCX Website: https://ccx.severalnines.com/
- Documentation: https://docs.severalnines.com/docs/ccx/
- Support: https://support.severalnines.com/hc/en-us/requests/new (sign up before you create the request).

## Requirement

- Terraform 0.13.x or later

## Quick Start

1. Sign up for CCX at https://ccx.severalnines.com/
2. Create a Terraform file called datastore.tf with the content below.
3. Set `client_id`, `client_secret` and datastore information. You can generate these credentials on the Account
   page (https://app.mydbservice.net/account) Authorization tab.

```terraform
terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 0.2.3"
    }
  }
}

provider "ccx" {
    client_id = "please_enter_your_client_id_here"
    client_secret = "please_enter_your_client_secret_here"
}

resource "ccx_datastore" "luna_postgres" {
  name           = "luna"
  size           = 1
  db_vendor      = "postgres"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
  network_type   = "public"
}

resource "ccx_datastore" "luna_mysql" {
  name           = "luna_mysql"
  size           = 3
  type           = "replication"
  db_vendor      = "mysql"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
  network_type   = "public"
}
```

Optionally you can create a VPC (supported by AWS)

```terraform
resource "ccx_vpc" "venus" {
    vpc_name = "venus"
    vpc_cloud_provider = "aws"
    vpc_cloud_region = "eu-north-1"
    vpc_ipv4_cidr = "10.10.0.0/16"
}
```

In that case set:

```terraform
    network_type = "private"
    network_vpc_uuid = ccx_vpc.venus.id
```

in the resource "ccx_datastore" section, see also [example_datastore.tf](examples/example_datastore.tf)

4. Run:

- `terraform init`
- `terraform apply `

5. Login to CCX and watch the datastore being deployed :)

## Advanced Usage

### Database Parameters

Database parameters can be configured for the cluster by using the block `db_params` inside the `ccx_datastore` block as follows:

```terraform
db_params {
   sql_mode = "STRICT"
}
```

The parameters will depend on the database vendor and version.
Refer to the `Settings > DB Parameters` section for a list of available parameters.

### Firewall Settings

Firewall settings can be configured for the cluster by using the block `firewall` inside the `ccx_datastore` block as follows:

```terraform
firewall {
   source = "x.x.x.x/32"
   description = "description here"
}

firewall {
   source = "y.y.y.y/32"
   description = "description here"
}
```

You may add multiple firewall blocks to allow multiple IP addresses.

### Notifications

Notifications can be configured for the cluster by including the following blocks inside the `ccx_datastore` block:

```terraform
notifications_enabled = true # or false
 notifications_emails = ["your@email.com", "your2@email.com"] # 
```

### Maintenance Settings

Maintenance settings can be configured for the cluster by including the following blocks inside the `ccx_datastore` block:

```terraform
maintenance_day_of_week = 2 # 1-7, 1 is Monday
maintenance_start_hour = 2 # 0-23
maintenance_end_hour = 4
```

### Scaling the cluster

Scaling the cluster can be done by changing the `size` parameter in the `ccx_datastore` block. When downscaling, the oldest non-primary node will be removed.

## Installing the provider from source

### Requirements

- go 1.21 or later
- Terraform 0.13.x or later

### Unix

Clone and use `make` to install.

1. Clone:   `git clone https://github.com/severalnines/terraform-provider-ccx`
2. Install: `make install`

### Windows

Clone, build and place the plugin in the right folder.

1. Clone: `git clone https://github.com/severalnines/terraform-provider-ccx`
2. Build: `go build -o ./bin/terraform-provider-ccx.exe ./cmd/terraform-provider-ccx`
3. Place: `./bin/terraform-provider-ccx.exe`
   as `%APPDATA%/terraform.d/plugins/registry.terraform.io/severalnines/ccx/0.2.0/windows_amd64/terraform-provider-ccx.exe`

## Using the provider

Create a provider and a resource file and specify account settings and datastore properties. The provider and resource
sections may be located in one file,
see https://github.com/severalnines/terraform-provider-ccx/examples/example_datastore.tf

### Create a terraform provider configuration

```
terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 0.2.3"
    }
  }
}
```

### Create a terraform provider file

```
provider  "ccx" {
	client_id  =  "your_ccx_client_id"
	client_secret  =  "your_ccx_client_secret
	# base_url = "optionally_use_a_different_base_url"
	# timeout = "15m"
}
```

> **Note:**
> 
> the option `base_url` may be used to specify a different ccx compliant cloud service.
> 
> The option `timeout` may be used to specify a different timeout for operations.
> Default is `"15m"`.
> Format is according to [ParseDuration](https://pkg.go.dev/time#ParseDuration).
> 

### Create a terraform resource file

```
resource "ccx_datastore" "luna" {
  name           = "luna"
  size           = 2
  db_vendor      = "postgres"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
  network_type   = "public"
}
```
or for mysql, replication (primary and two replicas, i.e size=3)
```
resource "ccx_datastore" "luna_mysql" {
  name           = "luna_mysql"
  size           = 3
  type           = "replication"
  db_vendor      = "mysql"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
  network_type   = "public"
}
```
### Create VPC

```
resource "ccx_vpc" "venus" {
  vpc_name = "venus"
  vpc_cloud_provider = "aws"
  vpc_cloud_region = "eu-north-1"
  vpc_ipv4_cidr = "10.10.0.0/16"
}
```

### Apply the created file

`terraform apply`

## Issues

If you have issues, please report them under Issues.
