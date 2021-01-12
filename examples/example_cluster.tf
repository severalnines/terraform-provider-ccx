provider "ccx" {
    address = "https://auth-api.s9s-dev.net/login"
    username = "simon+ccx@s9s.io"
    password = "Severalnines141$?"
}
resource "ccx_cluster" "spaceforce" {
    cluster_name = "spaceforce"
    cluster_type = "galera"
    cloud_provider = "aws"
    region = "eu-west-2"
    db_vendor = "mariadb"
    instance_size = "t3.medium"
    instance_iops = "100"
    db_username = "milen"
    db_password = "hristov"
    db_host = "10.11.12.13"
}