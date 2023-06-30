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
3. You must set client_id, client_secret and datastore information.

```
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
resource "ccx_datastore" "spaceforce" {
    name = "spaceforce"
    size = 1
    db_vendor = "mariadb"
    tags = ["new", "test"]
    cloud_provider = "aws"
    region = "eu-north-1"
    instance_size = "tiny"
    volume_size = 80
    volume_type = "gp2"
    network_type = "public"
}
```

Optionally you can create a VPC (supported by AWS)

```
resource "ccx_vpc" "venus" {
    vpc_name = "venus"
    vpc_cloud_provider = "aws"
    vpc_cloud_region = "eu-north-1"
    vpc_ipv4_cidr = "10.10.0.0/16"
}
```

In that case set:

```
    network_type = "private"
    network_vpc_uuid = ccx_vpc.venus.id
```

in the resource "ccx_datastore" section, see also [example_datastore.tf](examples/example_datastore.tf)

3. Run:

- `terraform init`
- `terraform apply `

4. Login to CCX and watch the datastore being deployed :)

## Installing the provider from source

### Requirements

- go 1.20 or later
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
}
```

> **Note:**
> the option `base_url` may be used to specify a different ccx compliant cloud service.

### Create a terraform resource file

```
resource "ccx_datastore" "luna" {
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
