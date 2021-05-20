terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 0.0.1"
    }
  }
}
provider "ccx" {
    auth_service_url = "https://ccx.s9s-dev.net/api/auth"
    username = "please_enter_your_username_here"
    password = "please_enter_your_password_here"
}
resource "ccx_cluster" "spaceforce" {
    cluster_name = "spaceforce"
    cluster_size = 1
    db_vendor = "mariadb"
    tags = ["new", "test"]
    cloud_provider = "aws"
    region = "eu-north-1"
    instance_size = "tiny"
    volume_size = 40
    volume_type = "gp2"
    network_type = "public"
}
