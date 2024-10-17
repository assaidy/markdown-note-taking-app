package main

import (
	"log"

	"github.com/assaidy/markdown-note-takin-app/server"
)

func main() {
	server := server.NewFiberServer()
	server.RegisterRoutes()
	log.Fatal(server.Listen(":8080"))
}
