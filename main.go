package main

import (
	"github.com/Courtcircuits/mitter-server/api"
	"github.com/Courtcircuits/mitter-server/util"
)

func main() {

	port := util.Get("PORT")

	server := api.NewServer(":" + port)

	err := server.Start()
	if err != nil {
		panic(err)
	}
}
