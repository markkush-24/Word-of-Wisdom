package app

import (
	"fmt"
	"word_of_wisdom/internal/client/network"
	logger "word_of_wisdom/internal/pkg/logging"
	"word_of_wisdom/internal/pkg/pow"
)

type ClientService struct {
	address string
	pow     pow.PowSolver
	conn    network.ConnectionService
	logger  logger.Logger
}

func NewClientService(
	address string,
	pow pow.PowSolver,
	connService network.ConnectionService,
	log logger.Logger) *ClientService {
	return &ClientService{
		address: address,
		pow:     pow,
		conn:    connService,
		logger:  log,
	}
}

func (c *ClientService) Start() error {
	conn, err := c.conn.Connect(c.address)
	if err != nil {
		c.logger.Error("Error connecting to the server", err)
		return fmt.Errorf("client service: failed to connect to server at %s: %w", c.address, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			c.logger.Error(fmt.Sprintf("Failed to close connection to %s: %v", c.address, err))
		}
	}()

	c.logger.Info("Connection to the server established.")

	// Read the challenge from the server.
	challenge, err := c.conn.ReadMessage(conn)
	if err != nil {
		c.logger.Error("Error reading challenge", err)
		return fmt.Errorf("client service: failed to read message: %w", err)
	}
	c.logger.Info(fmt.Sprintf("Challenge received: %s", challenge))

	// Generate PoW solution using the challenge.
	solution := c.pow.SolvePoW(challenge)
	c.logger.Info(fmt.Sprintf("PoW solution generated: %s", solution))

	// Send the solution to the server.
	err = c.conn.SendMessage(conn, solution)
	if err != nil {
		c.logger.Error("Error sending solution", err)
		return fmt.Errorf("client service: failed to send PoW solution: %w", err)
	}
	c.logger.Info("Solution sent to the server.")

	// Read the quote from the server.
	quote, err := c.conn.ReadMessage(conn)
	if err != nil {
		c.logger.Error("Error receiving quote", err)
		return fmt.Errorf("client service: failed to receive quote: %w", err)
	}
	c.logger.Info(fmt.Sprintf("Quote received: %s", quote))

	return nil
}
