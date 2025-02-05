package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	core "word_of_wisdom/config"
	"word_of_wisdom/internal/pkg/logging/zap"
	"word_of_wisdom/internal/pkg/quotes/repo"
	"word_of_wisdom/internal/server"
	"word_of_wisdom/internal/server/tcpserver"
	"word_of_wisdom/internal/server/verifier"
)

func main() {
	cfg := core.LoggerConfig{Level: "info"}
	log := zap.NewZapLogger(cfg)
	log.Info("Server starting")

	// Create a context to manage the application lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive system signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a signal handler for graceful shutdown
	go server.HandleShutdown(ctx, cancel, signalChan, log)

	// Initialize the quote repository
	quoteRepo, err := repo.NewQuoteRepository(repo.FileRepo)
	if err != nil {
		log.Error("Error initializing quote repository", err)
		// If an error occurs, cancel() can be called for immediate shutdown
		cancel()
		return
	}

	// Create a PoW verifier
	powVerifier := verifier.NewVerifier()

	// Create a PoW handler
	handler := &tcpserver.PoWHandler{
		PoWTarget:   core.PoWTarget,
		QuoteRepo:   quoteRepo,
		PowVerifier: powVerifier,
		Logger:      log,
	}

	// Create a WaitGroup to track server goroutines
	var wg sync.WaitGroup

	// Start the server in a separate goroutine and register it in the WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tcpserver.StartServer(ctx, core.ServerPort, handler, log)
		if err != nil {
			log.Error("Error starting server", err)
		}
	}()

	// The main goroutine waits for the context to be canceled (e.g., signal received)
	<-ctx.Done()
	log.Info("Server shutting down gracefully...")

	// Wait for all goroutines to finish (server loop, connection handlers, etc.)
	wg.Wait()
	log.Info("Server shutdown complete")
}
