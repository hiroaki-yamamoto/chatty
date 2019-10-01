package config

import (
	"context"
	"time"
)

// DBConfig represents a database configuration
type DBConfig struct {
	URI     string
	Name    string
	Timeout time.Duration
}

// TimeoutContext returns a context set a timeout with DBConfig.Timeout.
func (me DBConfig) TimeoutContext(ctx context.Context) (
	context.Context, context.CancelFunc,
) {
	return context.WithTimeout(ctx, me.Timeout)
}

// Config represents a configuration.
type Config struct {
	Db DBConfig
}
