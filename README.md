
# terraform-provider-ccx
**Requirements**
- go 1.14 or later
- Terraform 0.13.x or later

**Installing the provider**
--
 - Clone the master branch of the current repository
 `git clone https://github.com/severalnines/terraform-provider-ccx`
 - Build the terraform provider
- If using terraform < 0.14.0
`go build -o terraform-provider-ccx .`

`mkdir -p ~/.terraform.d/plugins/ && cp terraform-provider-ccx ~/.terraform.d/plugins/`
- If using terraform > 0.14.0 then the following steps are required:
- Step 1: Create the directory required for setup: `mkdir -p examples/.terraform.d/plugins/registry.terraform.io/hashicorp/ccx/1.1.0/linux_amd64/`
- Step 2: Execute the following command(s): `go build -o examples/.terraform.d/plugins/registry.terraform.io/hashicorp/ccx/0.1.0/linux_amd64/terraform-provider-ccx_v0.1.0 && rm -f examples/.terraform.lock.hcl && cd examples/  && terraform init -plugin-dir .terraform.d/plugins`
This will build the provider and place it in the correct directory. The provider will be available under the directory tree of examples only. If you wish to use the provider globaly , replace examples/ with your home dir (~).
`mkdir -p ~/.terraform.d/plugins/severalnines/ccx/1.0/amd64/ && cp terraform-provider-ccx ~/.terraform.d/plugins/severalnines/ccx/1.0/amd64/`

## **Using the provider**
- Create an terraform resource and provider file
```
provider  "ccx" {
	auth_service_url  =  "https://auth-api.s9s.io" 
	username  =  "your_ccx_email"
	password  =  "your_ccx_password"
}
```
- Create an terraform resource file
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
- Apply the created file
`terraform apply`

Bear in mind that this is experimental version and there might be some bugs
