
# terraform-provider-ccx
## Requirement
- go 1.14 or later
- Terraform 0.13.x or later

## Installing the provider
 - Clone the master branch of the current repository
 `git clone https://github.com/severalnines/terraform-provider-ccx`
### If using terraform <= 0.14.0

 - Build the terraform provider
`go build -o terraform-provider-ccx .`

`mkdir -p ~/.terraform.d/plugins/ && cp terraform-provider-ccx ~/.terraform.d/plugins/`

### If using terraform > 0.14.0
1. Create the directory required for setup: 
- `mkdir -p examples/.terraform.d/plugins/registry.terraform.io/hashicorp/ccx/1.1.0/linux_amd64/`
4. Execute the following command(s): 
- `go build -o examples/.terraform.d/plugins/registry.terraform.io/hashicorp/ccx/0.1.0/linux_amd64/terraform-provider-ccx_v0.1.0 && rm -f examples/.terraform.lock.hcl && cd examples/  && terraform init -plugin-dir .terraform.d/plugins`
This will build the provider and place it in the correct directory. The provider will be available under the directory tree of examples only. If you wish to use the provider globaly , replace examples/ with your home dir (~).
- `mkdir -p ~/.terraform.d/plugins/severalnines/ccx/1.0/amd64/ && cp terraform-provider-ccx ~/.terraform.d/plugins/severalnines/ccx/1.0/amd64/`

## Using the provider

Create a provider and a resource file and specify account settings and cluster properties. The provider and resource sections may be located in one file, see https://github.com/severalnines/terraform-provider-ccx/examples/example_cluster.tf

### Create an terraform provider file
```
provider  "ccx" {
	auth_service_url  =  "https://auth-api.s9s.io" 
	username  =  "your_ccx_email"
	password  =  "your_ccx_password"
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
    volume_size = 40
    volume_type = "gp2"
    network_type = "public"
}
```
### Apply the created file
`terraform apply`

## Issues
If you have issues, please report them under Issues.
