package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/severalnines/terraform-provider-ccx/terraform"
	"github.com/severalnines/terraform-provider-ccx/terraform/cluster"
	"github.com/severalnines/terraform-provider-ccx/terraform/vpc"
)

func main() {
	err := providerserver.Serve(
		context.Background(),
		func() provider.Provider {
			return terraform.New(&cluster.Resource{}, &vpc.Resource{})
		},
		providerserver.ServeOpts{
			Address: "severalnines/ccx/provider",
			// Debug:   true,
		},
	)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
