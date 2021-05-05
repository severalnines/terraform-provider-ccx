provider "ccx" {
    auth_service_url = "https://auth-api.s9s-dev.net"
    username = "simon+ccx@s9s.io"
    password = "Severalnines141$?"
}
resource "ccx_cluster" "spaceforce" {
    cluster_name = "spaceforce"
    cluster_size = 1
    db_vendor = "mariadb"
    tags = ["new", "test"]
    cloud_provider = "aws"
    region = "eu-west-2"
    instance_size = "t3.medium"
    volume_iops = 100
    volume_size = 40
    volume_type = "gp2"
    network_type = "public"
    network_ha_enabled = false
    network_vpc_uuid = "231321321321"
    network_az = ["a","b"]
}