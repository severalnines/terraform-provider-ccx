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
3. You must set client_id, client_secret and cluster information.
```
terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 1.5.0"
    }
  }
}
provider "ccx" {
    client_id = "please_enter_your_client_id_here"
    client_secret = "please_enter_your_client_secret_here"
}
resource "ccx_cluster" "spaceforce" {
    cluster_name = "spaceforce"
    cluster_size = 1
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
Optionally you can create a VPC (supported by GCP, AWS)
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
in the resource "ccx_cluster" section, see also [example_cluster_vpc.tf](examples/example_cluster.tf)

3. Run:
  - `terraform init`
  - `terraform apply `

4. Login to CCX and watch the cluster being deployed :)

## Installing the provider from source
### Requirements
- go 1.18 or later
- Terraform 0.13.x or later

### Unix
Clone and use `make` to install.
1. Clone:   `git clone https://github.com/severalnines/terraform-provider-ccx`
2. Install: `make install`

### Windows
Clone, build and place the plugin in the right folder.
1. Clone: `git clone https://github.com/severalnines/terraform-provider-ccx`
2. Build: `go build -o ./bin/terraform-provider-ccx.exe ./cmd/terraform-provider-ccx`
3. Place: `./bin/terraform-provider-ccx.exe` as `%APPDATA%/terraform.d/plugins/registry.terraform.io/severalnines/ccx/1.5.0/windows_amd64/terraform-provider-ccx.exe`

## Using the provider

Create a provider and a resource file and specify account settings and cluster properties. The provider and resource sections may be located in one file, see https://github.com/severalnines/terraform-provider-ccx/examples/example_cluster.tf
### Create a terraform provider configuration
```
terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 1.5.0"
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
resource "ccx_cluster" "spaceforce" {
  cluster_name = "spaceforce"
  cluster_size = 1
  db_vendor = "mariadb"
  db_version = "10.6"
  tags = ["new", "test"]
  cloud_provider = "aws"
  cloud_region = "eu-north-1"
  instance_size = "t3.medium"
  volume_size = 8000
  volume_type = "gp2"
  volume_iops = 0
  network_type = "public"
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

### Cluster with VPC peering (Optional)
> You can create a vpc and use it for peering

```
resource "ccx_cluster" "luna" {
  cluster_name = "luna"
  cluster_size = 1
  db_vendor = "mariadb"
  db_version = "10.6"
  tags = ["new", "test"]
  cloud_provider = "aws"
  cloud_region = "eu-north-1"
  instance_size = "t3.medium"
  volume_size = 9000
  volume_type = "gp2"
  volume_iops = 0
  network_type = "public"
}

resource "ccx_vpc" "venus" {
  vpc_name = "venus"
  vpc_cloud_provider = "aws"
  vpc_cloud_region = "eu-north-1"
  vpc_ipv4_cidr = "10.10.0.0/16"
}
```
Resource can be referenced in the network_vpc_uuid field in the following format:
```
<resource_type>.<resource_name>.<id>
```
## Issues
If you have issues, please report them under Issues.
