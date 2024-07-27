package main

import (
	"bufio"
	"fmt"
	"gamma/src/client"
	"log"
	"os"
	"strings"
)

func main() {
	c := client.NewClient("localhost:8080")
	err := c.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer c.Close()

	fmt.Println("Connected to server.")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Select service (AI or VAULT), or type 'exit' to quit: ")
		service, _ := reader.ReadString('\n')
		service = strings.TrimSpace(strings.ToUpper(service))

		if service == "EXIT" {
			break
		}

		if service != "AI" && service != "VAULT" {
			fmt.Println("Invalid service. Please choose AI or VAULT.")
			continue
		}

		err = c.SelectService(service)
		if err != nil {
			fmt.Printf("Error selecting service: %v\n", err)
			continue
		}

		fmt.Printf("Connected to %s service. Type your commands (type 'back' to change service, 'exit' to quit):\n", service)
		if service == "VAULT" {
			fmt.Println("Available commands: 'STORE <text>' to store text, 'RETRIEVE' to retrieve stored text")
		}

		for {
			fmt.Print("> ")
			command, _ := reader.ReadString('\n')
			command = strings.TrimSpace(command)

			if command == "back" {
				break
			}
			if command == "exit" {
				return
			}

			response, err := c.SendCommand(command)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("Server response:", response)
		}
	}
}
