provider "ccx" {
    username = "bob@severalnines.com"
    password = "H0hhw51@"
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

output "MOTD" {
  value = <<EOF
  ### Congratulations, your cluster ${ccx_cluster.spaceforce.cluster_name} with id ${ccx_cluster.spaceforce.id} 
  has been sucessfully deployed ### 
  ### Please visit: https://ccx.s9s-dev.net to view the status of this deployment
  EOF
}
