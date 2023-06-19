package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/severalnines/terraform-provider-ccx/terraform"
	"github.com/severalnines/terraform-provider-ccx/terraform/cluster"
	"github.com/severalnines/terraform-provider-ccx/terraform/vpc"
)

func main() {
	p := terraform.New(
		&cluster.Resource{},
		&vpc.Resource{},
	)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: p.Resources,
	})
}
