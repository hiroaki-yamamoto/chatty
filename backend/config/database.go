package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB represents a database configuration
type DB struct {
	Timeout time.Duration
	URI     string
	Name    string
}

// TimeoutContext returns a context set a timeout with DB.Timeout.
func (me DB) TimeoutContext(ctx context.Context) (
	context.Context, context.CancelFunc,
) {
	return context.WithTimeout(ctx, me.Timeout)
}

// CreateClient creates a client to connect to the DB.
func (me DB) CreateClient() (*mongo.Client, error) {
	return mongo.NewClient(
		options.Client().ApplyURI(me.URI),
	)
}
