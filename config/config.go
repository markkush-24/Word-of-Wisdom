package core

const (
	ServerAddress       = "localhost:8080"
	ServerAddressDocker = "server-container:8080"
	PoWTarget           = 5

	ConnectionType      = "tcp"
	ServerPort          = 8080
	HashPrefixCharacter = "0"               // Character that the hash should start with
	CheckpointInterval  = 100000            // Interval for checkpoints to log progress
	SolutionPrefix      = "solutionNumber:" // Prefix for the solution
)

const (
	ConnectionTypeTCP = "tcp"
	ConnectionTypeUDP = "udp"
)

// LoggerConfig for logger configuration
type LoggerConfig struct {
	DisableCaller     bool   `env:"CALLER,default=false"`
	DisableStacktrace bool   `env:"STACKTRACE,default=false"`
	Level             string `env:"LEVEL,default=debug"`
}

// AppConfig for application configuration
type AppConfig struct {
	QuoteSource    string `env:"QUOTE_SOURCE,default=json"` // "file" or "json"
	ConnectionType string `env:"CONNECTION_TYPE,default=tcp"`
}
