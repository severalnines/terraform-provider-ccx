package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/severalnines/terraform-provider-ccx/resources"
)

func main() {
	p := resources.Provider(
		&resources.Datastore{},
		&resources.VPC{},
	)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: p.Resources,
	})
}
