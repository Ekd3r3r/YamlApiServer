package main

import (
	"YamlApiServer/pkg/server"
)

func main() {
	server := server.NewServer()
	server.Run(":8080")
}
