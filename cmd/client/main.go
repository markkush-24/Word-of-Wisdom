package main

import (
	core "word_of_wisdom/config"
	"word_of_wisdom/internal/client/app"
	"word_of_wisdom/internal/client/network"
	"word_of_wisdom/internal/pkg/logging/zap"
	"word_of_wisdom/internal/pkg/pow/powsolver"
)

func main() {
	// Initialize logger
	cfg := core.LoggerConfig{Level: "info"}
	log := zap.NewZapLogger(cfg)
	log.Info("Client application started")

	// Initialize application configuration
	appConfig := core.AppConfig{
		ConnectionType: core.ConnectionTypeTCP,
	}

	// Get connection service via the factory method
	connService, err := network.ConnectionFactory(appConfig)
	if err != nil {
		log.Error("Error creating connection service", err)
		log.Fatalf("Error: %v", err)
	}

	// Pass all necessary arguments to NewClientService.
	// Note: The PoW solver now expects a challenge.
	clientService := app.NewClientService(
		core.ServerAddressDocker, powsolver.NewSimplePoW(core.PoWTarget, log), connService, log)

	// Start the client.
	err = clientService.Start()
	if err != nil {
		log.Error("Error starting the client", err)
		log.Fatalf("Error: %v", err)
	}
}
