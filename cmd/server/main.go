package main

import (
	"gamma/src/server"
	"log"
)

func main() {
	s, err := server.NewServer(":8080", ":8081")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Println("Starting server...")
	err = s.Start()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
