package config

// Config represents a configuration.
type Config struct {
	Db        DB
	Servers   map[string]*Server
	Broker    Broker
	Recaptcha string
}
