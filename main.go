package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/severalnines/terraform-provider-ccx/terraform"
	"github.com/severalnines/terraform-provider-ccx/terraform/datastore"
	"github.com/severalnines/terraform-provider-ccx/terraform/vpc"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	p := terraform.New(
		&datastore.Resource{},
		&vpc.Resource{},
	)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: p.Resources,
	})
}
