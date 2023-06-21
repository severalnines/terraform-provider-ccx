terraform {
    required_providers {
        ccx = {
            source  = "severalnines/ccx"
            version = "~> 1.5.0"
        }
    }
}

provider "ccx" {
    client_id = "replace_with_your_own_client_id"
    client_secret = "replace_with_your_own_client_secret"
    # base_url = "optionally_use_a_different_base_url"
}

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
    ipv4_cidr = "10.10.0.0/16"
}

output "MOTD" {
  value = <<EOF
  ### Congratulations, your cluster ${ccx_cluster.spaceforce.cluster_name} with id ${ccx_cluster.spaceforce.id} 
  has been sucessfully deployed ### 
  ### Please visit: https://ccx.s9s-dev.net to view the status of this deployment
  EOF
}

