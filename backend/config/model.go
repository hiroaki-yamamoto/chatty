package config

// Config represents a configuration.
type Config struct {
	Db        DB
	Server    Server
	Broker    Broker
	Recaptcha string
}
