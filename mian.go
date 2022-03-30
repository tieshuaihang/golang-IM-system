package main

import "golangIM/server"

func main() {
	server := server.NewServer("0.0.0.0", 8888)
	server.Start()
}
