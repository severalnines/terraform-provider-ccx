package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/severalnines/terraform-provider-ccx/resources"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	p := resources.Provider(
		&resources.Datastore{},
		&resources.VPC{},
	)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: p.Resources,
	})
}
