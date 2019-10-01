package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// DB represents a database configuration
type DB struct {
	Timeout  time.Duration
	Client   *mongo.Client
	Database *mongo.Database
}

// TimeoutContext returns a context set a timeout with DBConfig.Timeout.
func (me DB) TimeoutContext(ctx context.Context) (
	context.Context, context.CancelFunc,
) {
	return context.WithTimeout(ctx, me.Timeout)
}

// Config represents a configuration.
type Config struct {
	Db DB
}
