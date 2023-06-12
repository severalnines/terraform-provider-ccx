terraform {
    required_providers {
        ccx = {
            source  = "severalnines/ccx"
            version = "~> 1.1.0"
        }
    }
}
provider "ccx" {
    client_id = "_" # client_id from ccx
    client_secret = "_" # client_secret from ccx

    # auth_endpoint = "" # uncomment to specift custom endpoint for auth API calls
    # default value = "https://app.mydbservice.net/api/auth"

    # cluster_endpoint = "" # uncomment to specift custom endpoint for cluster API calls
    # default value = "https://app.mydbservice.net/api/prov/api/v2/cluster"

    # vpc_endpoint = "" # uncomment to specift custom endpoint for vpc API calls
    # default value = "https://app.mydbservice.net/api/prov/api/v2/cluster"


    # dev_mode = true # uncomment to enable dev_mode; does not call any API; uses in memory fake api

    # mock_file = "mock.json" # uncomment to load mock state
    # path to mock file, check .mock.sample.json for an example
    # can be used in conjunction with .terraform.tfstate for manual testing
}
resource "ccx_cluster" "spaceforce" {
    cluster_name = "spaceforce"
    cluster_size = 1
    db_vendor = "mariadb"
    tags = ["new", "test"]
#    cloud_space    = ""
    cloud_provider = "aws"
    cloud_region = "eu-north-1"
    instance_size = "tiny"
    volume_size = 80
    volume_type = "gp2"
    network_type = "public"
}

resource "ccx_cluster" "luna" {
    cluster_name = "luna"
    cluster_size = 1
    db_vendor = "mariadb"
    tags = ["new", "test"]
    cloud_provider = "aws"
    cloud_region = "eu-north-1"
    instance_size = "tiny"
    volume_size = 80
    volume_type = "gp2"
    network_type = "public"
}

resource "ccx_vpc" "venus" {
    name = "venus"
    cloud_provider = "aws"
    cloud_region = "eu-north-1"
}