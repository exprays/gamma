package main

import (
	"log"
	"gamma/src/server"
)

func main() {
	s := server.NewServer(":8080")
	log.Println("Starting server on :8080")
	err := s.Start()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
