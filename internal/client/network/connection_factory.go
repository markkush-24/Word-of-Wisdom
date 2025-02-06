package network

import (
	"word_of_wisdom/config" //nolint:goimports
)

// ConnectionFactory factory for creating connection services
func ConnectionFactory(cfg core.AppConfig) (ConnectionService, error) {
	return GetConnectionService(cfg.ConnectionType)
}
