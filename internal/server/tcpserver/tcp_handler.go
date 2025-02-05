package tcpserver

import (
	"context"
	"fmt"
	"net"
	logger "word_of_wisdom/internal/pkg/logging"
	"word_of_wisdom/internal/pkg/quotes/repo"
	"word_of_wisdom/internal/server/network"
	"word_of_wisdom/internal/server/verifier"
)

type TCPHandler interface {
	Handle(ctx context.Context, conn net.Conn)
}

type PoWHandler struct {
	PoWTarget   int
	QuoteRepo   repo.QuoteRepo
	PowVerifier verifier.PoWVerifier
	Logger      logger.Logger
}

func (h *PoWHandler) Handle(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			h.Logger.Error(fmt.Sprintf("Error closing connection: %v", err))
		}
	}()

	// Processing PoW
	solution, err := network.HandlePoWRequest(conn, h.Logger)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error processing PoW: %v", err))
		return
	}

	// Verifying PoW
	if !h.PowVerifier.VerifyPoW(solution, h.PoWTarget) {
		h.Logger.Info("Incorrect PoW solution. Requesting a new solution")
		if _, err := fmt.Fprintf(conn, "Incorrect PoW solution. Please try again.\n"); err != nil {
			h.Logger.Error(fmt.Sprintf("Error sending incorrect PoW message to client: %v", err))
		}
		return
	}

	// Retrieving quote
	quote, err := h.QuoteRepo.GetQuote(ctx)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error retrieving quote: %v", err))
		return
	}

	// Sending quote to the client
	if err := network.SendQuoteToClient(conn, quote, h.Logger); err != nil {
		h.Logger.Error(fmt.Sprintf("Error sending quote to client: %v", err))
	}
}
