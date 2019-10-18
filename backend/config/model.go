package config

import (
	"context"
	"time"
)

// Config represents a configuration.
type Config struct {
	Timeout time.Duration
	Db      DB
	Server  Server
	Broker  Broker
}

// TimeoutContext returns a context set a timeout with DB.Timeout.
func (me Config) TimeoutContext(ctx context.Context) (
	context.Context, context.CancelFunc,
) {
	return context.WithTimeout(ctx, me.Timeout)
}
