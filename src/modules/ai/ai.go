package ai

import "fmt"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ProcessCommand(command string) string {
	// TODO: Implement actual AI processing here
	return fmt.Sprintf("Received command: %s. This is a dummy response.", command)
}
