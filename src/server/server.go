package server

import (
	"bufio"
	"fmt"
	"gamma/src/modules/ai"
	"gamma/src/modules/vault"
	"net"
	"strings"
)

type Server struct {
	addr         string
	aiService    *ai.Service
	vaultService *vault.Service
}

func NewServer(addr string) *Server {
	return &Server{
		addr:         addr,
		aiService:    ai.NewService(),
		vaultService: vault.NewService("vault.txt"),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		serviceType, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading service type: %v\n", err)
			return
		}
		serviceType = strings.TrimSpace(strings.ToUpper(serviceType))

		var service interface {
			ProcessCommand(string) string
		}

		switch serviceType {
		case "AI":
			service = s.aiService
		case "VAULT":
			service = s.vaultService
		default:
			fmt.Fprintf(conn, "Error: Invalid service type\n")
			continue
		}

		fmt.Fprintf(conn, "Connected to %s service\n", serviceType)

		for {
			command, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading command: %v\n", err)
				return
			}
			command = strings.TrimSpace(command)

			if command == "back" {
				break
			}

			response := service.ProcessCommand(command)
			fmt.Fprintf(conn, "%s\n", response)
		}
	}
}
