package server

import (
	"bufio"
	"fmt"
	"gamma/src/modules/ai"
	"gamma/src/modules/vault"
	"net"
	"net/http"
	"strings"
)

type Server struct {
	tcpAddr      string
	httpAddr     string
	aiService    *ai.Service
	vaultService *vault.Service
}

func NewServer(tcpAddr, httpAddr string) (*Server, error) {
	vaultService, err := vault.NewService("vault.db", httpAddr)
	if err != nil {
		return nil, fmt.Errorf("error creating vault service: %v", err)
	}

	return &Server{
		tcpAddr:      tcpAddr,
		httpAddr:     httpAddr,
		aiService:    ai.NewService(),
		vaultService: vaultService,
	}, nil
}

func (s *Server) Start() error {
	// Start HTTP server for file uploads
	go func() {
		http.HandleFunc("/upload", s.vaultService.HandleFileUpload)
		fmt.Printf("Starting HTTP server on %s\n", s.httpAddr)
		if err := http.ListenAndServe(s.httpAddr, nil); err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
		}
	}()

	// Start TCP server
	listener, err := net.Listen("tcp", s.tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Starting TCP server on %s\n", s.tcpAddr)
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

	fmt.Fprintf(conn, "Connected to server.\n")
	fmt.Fprintf(conn, "Select service (AI or VAULT), or type 'exit' to quit: ")

	for {
		serviceType, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading service type: %v\n", err)
			return
		}
		serviceType = strings.TrimSpace(strings.ToUpper(serviceType))

		if serviceType == "EXIT" {
			fmt.Fprintf(conn, "Goodbye!\n")
			return
		}

		var service interface {
			ProcessCommand(string) string
		}

		switch serviceType {
		case "AI":
			service = s.aiService
		case "VAULT":
			service = s.vaultService
		default:
			fmt.Fprintf(conn, "Error: Invalid service type. Please choose AI or VAULT.\n")
			continue
		}

		fmt.Fprintf(conn, "Connected to %s service. Type your commands (type 'back' to change service, 'exit' to quit):\n", serviceType)

		for {
			fmt.Fprintf(conn, "> ")
			command, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading command: %v\n", err)
				return
			}
			command = strings.TrimSpace(command)

			if command == "back" {
				break
			}
			if command == "exit" {
				fmt.Fprintf(conn, "Goodbye!\n")
				return
			}

			response := service.ProcessCommand(command)
			fmt.Fprintf(conn, "Server response: %s\n", response)
		}

		fmt.Fprintf(conn, "Select service (AI or VAULT), or type 'exit' to quit: ")
	}
}
