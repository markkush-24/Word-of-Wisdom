package network

import (
	"fmt"
	"net"

	logger "word_of_wisdom/internal/pkg/logging"
)

type TCPConnectionService struct {
	logger logger.Logger
}

// Connect establishes a TCP connection to the specified address.
// It returns the connection or an error if the connection fails.
func (s *TCPConnectionService) Connect(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to connect to the server: %v", err))
		return nil, fmt.Errorf("TCPConnectionService: failed to dial %s: %w", address, err)
	}
	return conn, nil
}

// ReadMessage reads a message from the provided TCP connection.
// It returns the message as a string or an error if reading fails.
func (s *TCPConnectionService) ReadMessage(conn net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error reading message from server: %v", err))
		return "", fmt.Errorf("TCPConnectionService: failed to read message: %w", err)
	}
	return string(buffer[:n]), nil
}

// SendMessage sends the specified message through the provided TCP connection.
// It returns an error if the sending fails.
func (s *TCPConnectionService) SendMessage(conn net.Conn, message string) error {
	if _, err := fmt.Fprintf(conn, "%s\n", message); err != nil {
		s.logger.Error(fmt.Sprintf("Error sending message: %v", err))
		return fmt.Errorf("TCPConnectionService: failed to send message: %w", err)
	}
	return nil
}
