package main

import "fmt"

func main() {
	c := CCXAuth("milen", "password")
	id, token := c.GetUserId()
	print(id)
	d := GetClusters(id, token)
	fmt.Println(d)
}
