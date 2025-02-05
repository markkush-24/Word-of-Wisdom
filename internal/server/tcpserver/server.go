package tcpserver

import (
	"context"
	"fmt"
	"net"
	"strconv"
	core "word_of_wisdom/config"
	logger "word_of_wisdom/internal/pkg/logging"
)

// StartServer starts the server and handles connections using a pool.
func StartServer(ctx context.Context, port int, handler TCPHandler, log logger.Logger) error {
	listener, err := net.Listen(core.ConnectionType, ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Info("Server started", "port", port)

	// Close the listener when the work is completed
	defer func() {
		if err := listener.Close(); err != nil {
			log.Warn("Failed to close listener", "error", err)
		}
	}()

	connPool := make(chan net.Conn)

	// Goroutine for accepting connections
	go func() {
		defer close(connPool)
		for {
			conn, err := listener.Accept()
			if ctx.Err() != nil { // If shutdown signal received, exit
				if err == nil {
					closeErr := conn.Close()
					if closeErr != nil {
						log.Warn("Error closing connection", "error", closeErr)
					}
				}
				break
			}

			if err != nil {
				log.Warn("Error accepting connection", "error", err)
				continue
			}

			connPool <- conn
		}
	}()

	for {
		select {
		case <-ctx.Done(): // When the context is canceled, shutdown the server
			return ctx.Err()
		case conn, ok := <-connPool:
			if !ok {
				return ctx.Err()
			}

			// Pass the connection to be handled by a new handler
			log.Info("Connection accepted", "remote", conn.RemoteAddr())
			go handler.Handle(ctx, conn)
		}
	}
}
