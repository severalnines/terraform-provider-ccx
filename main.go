package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/severalnines/terraform-provider-ccx/resources"
)

func main() {
	p := resources.Provider(
		&resources.Datastore{},
		&resources.VPC{},
	)

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: "registry.terraform.io/severalnines/ccx",
		ProviderFunc: p.Resources,
	}

	plugin.Serve(opts)
}
