provider "ccx" {
    auth_service_url = "https://ccx.s9s-dev.net/api/auth"
    username = "simon+ccx@s9s.io"
    password = "Severalnines141$?"
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
    network_type = "private"
    network_vpc_uuid =ccx_vpc.newVpc.id
}

resource "ccx_vpc" "newVpc" {
    vpc_name = "spaceforce_vpc"
    vpc_cloud_provider = "aws"
    vpc_cloud_region = "eu-north-1"
    vpc_ipv4_cidr = "10.10.0.0/16"
}
output "MOTD" {
  value = <<EOF
  ### Congratulations, your cluster ${ccx_cluster.spaceforce.cluster_name} with id ${ccx_cluster.spaceforce.id} 
  has been sucessfully deployed ### 
  ### Please visit: https://ccx.s9s-dev.net to view the status of this deployment
  EOF
}