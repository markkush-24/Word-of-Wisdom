package network

import "fmt"

// ConnectionConstructor defines a constructor function for creating a ConnectionService.
type ConnectionConstructor func() ConnectionService

var connectionRegistry = make(map[string]ConnectionConstructor)

// RegisterConnectionType registers a new connection type with its constructor.
func RegisterConnectionType(name string, constructor ConnectionConstructor) {
	connectionRegistry[name] = constructor
}

// GetConnectionService returns a ConnectionService for the specified type.
func GetConnectionService(connectionType string) (ConnectionService, error) {
	constructor, ok := connectionRegistry[connectionType]
	if !ok {
		return nil, fmt.Errorf("unsupported connection type: %s", connectionType)
	}
	return constructor(), nil
}
