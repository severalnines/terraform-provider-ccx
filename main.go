package main

import "github.com/severalnines/terraform-provider-ccx/services"

func main() {
	c := services.CCXAuth("milen", "password")
	id, token := c.GetUserId()
	services.GetClusters(id, token)
	services.GetClusterByID("41def0d0-d447-4bf2-9990-9a5066bd8b88", token)
	services.CreateCluster(
		id, "spaceforce", "mariadb", "aws", "eu-west-2",
		"mariadb", "t3.medium", 300, "milen", "dd0953", "0.0.0.0/0", token,
	)

}
