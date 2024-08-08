terraform {
  required_providers {
    ccx = {
      source  = "severalnines/ccx"
      version = "~> 0.2.3"
    }
  }
}

provider "ccx" {
  client_id     = "replace_with_your_own_client_id"
  client_secret = "replace_with_your_own_client_secret"
  # base_url = "optionally_use_a_different_base_url"
}


resource "ccx_datastore" "luna" {
  name           = "luna"
  size           = 1
  db_vendor      = "postgres"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
  network_type   = "public"

  db_params = {
    deadlock_timeout = "2"
  }

  firewall {
    source = "2.3.41.5/10"
    description = "hello"
  }

  firewall {
    source = "2.2.2.2/24"
    description = "world"
  }
}

output "MOTD" {
  value = <<EOF
  ### Congratulations, your datastore ${ccx_datastore.luna.name} with id ${ccx_datastore.luna.id}
  has been sucessfully created ###
  ### Please visit: https://app.mydbservice.net/ to view the status of its deployment
  EOF
}

