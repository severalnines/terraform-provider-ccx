
# terraform-provider-ccx

This is the Terraform Provider for the Severalnines CCX Database as a service platform.
- CCX Website: https://ccx.severalnines.com/
- Documentation: https://docs.severalnines.com/docs/ccx/
- Support: https://support.severalnines.com/hc/en-us/requests/new (sign up before you create the request).

## Requirement
- Terraform 0.13.x or later

## Quick Start
1. Sign up for CCX at https://ccx.severalnines.com/
2. Create a Terraform file called datastore.tf with the content below. You must set client_id, client_secret and you may change cluster_name,db_vendor,region, and instance_size (tiny, small, medium, large).
```
terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 0.0.1"
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
resource "ccx_vpc" "newVpc" {
    vpc_name = "spaceforce_vpc"
    vpc_cloud_provider = "aws"
    vpc_cloud_region = "eu-north-1"
    vpc_ipv4_cidr = "10.10.0.0/16"
}
```
In that case set:
```
    network_type = "private"
    network_vpc_uuid =ccx_vpc.newVpc.id
```
in the resource "ccx_cluster" section, see also example_cluster_vpc.tf

3. Run:
  * `terraform init`
  * `terraform appy`

4. Login to CCX and watch the cluster being deployed :)

## Installing the provider from source
### Requirement
- go 1.14 or later
- Terraform 0.13.x or later

### Clone the master branch of the current repository
 - `git clone https://github.com/severalnines/terraform-provider-ccx`
 
### If using terraform <= 0.14.0

 - Build the terraform provider
`go build -o terraform-provider-ccx .`

`mkdir -p ~/.terraform.d/plugins/ && cp terraform-provider-ccx ~/.terraform.d/plugins/`

### If using terraform > 0.14.0
1. Create the directory required for setup: 
- `mkdir -p examples/.terraform.d/plugins/registry.terraform.io/hashicorp/ccx/1.1.0/linux_amd64/`
2. Execute the following command(s): 
- `go build -o examples/.terraform.d/plugins/registry.terraform.io/hashicorp/ccx/0.1.0/linux_amd64/terraform-provider-ccx_v0.1.0 && rm -f examples/.terraform.lock.hcl && cd examples/  && terraform init -plugin-dir .terraform.d/plugins`
This will build the provider and place it in the correct directory. The provider will be available under the directory tree of examples only. If you wish to use the provider globaly , replace examples/ with your home dir (~).


## Using the provider

Create a provider and a resource file and specify account settings and cluster properties. The provider and resource sections may be located in one file, see https://github.com/severalnines/terraform-provider-ccx/examples/example_cluster.tf
### Create a terraform provider configuration
```
terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 0.0.1"
    }
  }
}
```
### Create an terraform provider file
```
provider  "ccx" {
	client_id  =  "your_ccx_client_id"
	client_secret  =  "your_ccx_client_secret
}
```
### Create an terraform resource file
```
resource "ccx_cluster" "spaceforce" {
    cluster_name = "spaceforce"
    cluster_size = 1
    db_vendor = "mariadb"
    tags = ["new", "test"]
    cloud_provider = "aws"
    region = "eu-north-1"
    instance_size = "tiny"
    volume_iops = 100
    volume_size = 80
    volume_type = "gp2"
    network_type = "public"
}
```
### Create VPC ( Used for VPC Peering )
```
resource "ccx_vpc" "newVpc" {
    vpc_name = "spaceforce_vpc"
    vpc_cloud_provider = "aws"
    vpc_cloud_region = "eu-north-1"
    vpc_ipv4_cidr = "10.10.0.0/16"
}
```

### Apply the created file
`terraform apply`

### Optional: You can use the VPC Created above to  deploy a cluster
```
resource "ccx_cluster" "spaceforce" {
    cluster_name = "spaceforce"
    cluster_size = 1
    db_vendor = "mariadb"
    tags = ["new", "test"]
    cloud_provider = "aws"
    region = "eu-north-1"
    instance_size = "tiny"
    volume_iops = 100
    volume_size = 80
    volume_type = "gp2"
    network_type = "public"
    network_vpc_uuid = ccx_vpc.newVpc.id
}
```
Resource can be referenced in the network_vpc_uuid field in the following format:
```
<resource_type>.<resource_name>.<id>
```
## Issues
If you have issues, please report them under Issues.
