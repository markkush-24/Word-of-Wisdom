package network

import (
	"fmt"
	"net"
	"time"
	logger "word_of_wisdom/internal/pkg/logging"
	"word_of_wisdom/internal/pkg/quotes/model"
)

// HandlePoWRequest sends a unique challenge to the client and waits for the PoW solution.
func HandlePoWRequest(conn net.Conn, log logger.Logger) (string, error) {
	// Generate a unique challenge (e.g., based on the current time)
	challenge := fmt.Sprintf("challenge:%d", time.Now().UnixNano())
	log.Info("Sending PoW challenge to the client:", challenge)

	// Send the challenge to the client
	if _, err := fmt.Fprintf(conn, "%s\n", challenge); err != nil {
		log.Error(fmt.Sprintf("Error sending PoW challenge to the client: %v", err))
		return "", err
	}

	// Wait for the client to send the PoW solution
	var solution string
	log.Info("Waiting for PoW solution from the client")
	if _, err := fmt.Fscanf(conn, "%s", &solution); err != nil {
		log.Error(fmt.Sprintf("Error receiving PoW solution from the client: %v", err))
		return "", err
	}
	log.Info(fmt.Sprintf("PoW solution received from the client: %s", solution))

	return solution, nil
}

// SendQuoteToClient sends a quote and its author to the client
func SendQuoteToClient(conn net.Conn, quote model.Quote, log logger.Logger) error {
	quoteMessage := fmt.Sprintf("Quote: \"%s\" - Author: %s", quote.Quote, quote.Author)
	if _, err := fmt.Fprintf(conn, "%s", quoteMessage); err != nil {
		log.Error(fmt.Sprintf("Error sending quote to the client: %v", err))
		return err
	}
	log.Info(fmt.Sprintf("Quote sent to the client: %s", quoteMessage))
	return nil
}
