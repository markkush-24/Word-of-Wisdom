package network

import "net"

// ConnectionService interface for working with connections
type ConnectionService interface {
	Connect(address string) (net.Conn, error)
	ReadMessage(conn net.Conn) (string, error)
	SendMessage(conn net.Conn, message string) error
}
