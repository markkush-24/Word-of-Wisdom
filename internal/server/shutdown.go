package server

import (
	"context"
	"os"
	logger "word_of_wisdom/internal/pkg/logging"
)

// HandleShutdown listens for OS signals and initiates graceful shutdown
func HandleShutdown(ctx context.Context, cancel context.CancelFunc, signalChan <-chan os.Signal, log logger.Logger) {
	select {
	case <-signalChan:
		log.Info("Shutdown signal received. Shutting down the server...")
		cancel()
	case <-ctx.Done():
		log.Info("Shutdown already initiated.")
	}
}
