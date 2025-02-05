package network

import (
	"fmt"
	"word_of_wisdom/config" //nolint:goimports
)

// ConnectionFactory factory for creating connection services
func ConnectionFactory(config core.AppConfig) (ConnectionService, error) {
	switch core.ConnectionType {
	case core.ConnectionTypeTCP:
		return &TCPConnectionService{}, nil
	default:
		return nil, fmt.Errorf("unsupported connection type: %s", config.ConnectionType)
	}
}
