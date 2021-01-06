package main

func main() {
	c := CCXAuth("milen", "password")
	id, token := c.GetUserId()
	GetClusters(id, token)

}
