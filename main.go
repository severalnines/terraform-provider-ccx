package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/severalnines/terraform-provider-ccx/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
