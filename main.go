package main

func main() {
	c := CCXAuth("milen", "password")
	id, token := c.GetUserId()
	GetClusters(id, token)
	GetClusterByID("41def0d0-d447-4bf2-9990-9a5066bd8b88", token)

}
