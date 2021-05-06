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
    volume_iops = 100
    volume_size = 40
    volume_type = "gp2"
    network_type = "public"
}