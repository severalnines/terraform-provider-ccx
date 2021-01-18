
# terraform-provider-ccx

**Installing the provider**
--
 - Clone the master branch of the current repository
 `git clone https://github.com/severalnines/terraform-provider-ccx`
 - Build the terraform provider
 `go build -o terraform-provider-ccx .`
- If using terraform < 14.0
`mkdir -p ~/.terraform.d/plugins/ && cp terraform-provider-ccx ~/.terraform.d/plugins/`
- If using terraform > 14.0
`mkdir -p ~/.terraform.d/plugins/severalnines/ccx/1.0/amd64/ && cp terraform-provider-ccx ~/.terraform.d/plugins/severalnines/ccx/1.0/amd64/`

## **Using the provider**
- Create work directory
`mkdir -p myAwesomeCCXProject`
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
resource  "ccx_cluster"  "your_resource_name" {
	cluster_name  =  "your_cluster_name"
	cluster_type  =  "galera" 
	cloud_provider  =  "aws"
	region  =  "eu-west-2"
	db_vendor  =  "mariadb"
	instance_size  =  "t3.medium"
	instance_iops  =  100
	db_username  =  "your_prefered_username"
	db_password  =  "your_prefered_password"
	db_host  =  "host_from_which_connections_are_allowed"
}
```
- Apply the created file
`terraform apply`

Bear in mind that this is experimental version and there might be some bugs
