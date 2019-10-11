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

// CreateClient creates a client to connect to the DB.
func (me DB) CreateClient() (*mongo.Client, error) {
	return mongo.NewClient(
		options.Client().ApplyURI(me.URI),
	)
}

// TimeoutContext returns a context set a timeout with DB.Timeout.
func (me DB) TimeoutContext(ctx context.Context) (
	context.Context, context.CancelFunc,
) {
	return context.WithTimeout(ctx, me.Timeout)
}

// Server represents a sever configuration.
type Server struct {
	Type    string
	Addr    string
	Timeout time.Duration // Operation time limit
}

// TimeoutContext returns a context set a timeout with DB.Timeout.
func (me Server) TimeoutContext(ctx context.Context) (
	context.Context, context.CancelFunc,
) {
	return context.WithTimeout(ctx, me.Timeout)
}

// Config represents a configuration.
type Config struct {
	Db     DB
	Server Server
}
