package vault

import (
	"fmt"
	"os"
	"strings"
)

type Service struct {
	filePath string
}

func NewService(filePath string) *Service {
	return &Service{filePath: filePath}
}

func (s *Service) ProcessCommand(command string) string {
	parts := strings.SplitN(command, " ", 2)
	if len(parts) < 2 {
		return "Invalid command. Use 'STORE <text>' or 'RETRIEVE'."
	}

	switch parts[0] {
	case "STORE":
		return s.storeText(parts[1])
	case "RETRIEVE":
		return s.retrieveText()
	default:
		return "Unknown command. Use 'STORE <text>' or 'RETRIEVE'."
	}
}

func (s *Service) storeText(text string) string {
	err := os.WriteFile(s.filePath, []byte(text), 0644)
	if err != nil {
		return fmt.Sprintf("Error storing text: %v", err)
	}
	return "Text stored successfully."
}

func (s *Service) retrieveText() string {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "No text stored yet."
		}
		return fmt.Sprintf("Error retrieving text: %v", err)
	}
	return fmt.Sprintf("Retrieved text: %s", string(data))
}
