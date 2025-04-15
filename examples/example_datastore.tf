terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 0.4.6"
    }
  }
}

provider "ccx" {
  client_id     = "replace_with_your_own_client_id"
  client_secret = "replace_with_your_own_client_secret"
  # base_url = "optionally_use_a_different_base_url"
}


resource "ccx_datastore" "luna" {
  name            = "luna"
  size            = 1
  db_vendor       = "mariadb"
  tags            = ["new", "test"]
  cloud_provider  = "aws"
  cloud_region    = "eu-north-1"
  instance_size   = "m5.large"
  volume_size     = 80
  volume_type     = "gp2"
  parameter_group = ccx_parameter_group.asteroid.id

  firewall {
    source      = "2.3.41.5/32"
    description = "hello"
  }

  firewall {
    source      = "2.2.2.0/24"
    description = "world"
  }

  notifications_enabled = true # or false
  notifications_emails  = ["your@email.com"] # list of emails

  maintenance_day_of_week = 1 # 1-7, 1 is Monday
  maintenance_start_hour  = 2 # 0-23
  maintenance_end_hour    = 4
}

resource "ccx_parameter_group" "asteroid" {
    name             = "asteroid"
    database_vendor  = "mariadb"
    database_version = "10.11"
    database_type    = "galera"

    parameters = {
      table_open_cache = 8000
      sql_mode         = "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION"
    }
}

output "MOTD" {
  value = <<EOF
  ### Congratulations, your datastore ${ccx_datastore.luna.name} with id ${ccx_datastore.luna.id}
  has been sucessfully created ###
  ### Please visit: https://app.mydbservice.net/ to view the status of its deployment
  EOF
}

