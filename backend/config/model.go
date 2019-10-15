package config

import (
	"context"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WithTimeout represents a timeout structure
type WithTimeout struct {
	Timeout time.Duration
}

// TimeoutContext returns a context set a timeout with DB.Timeout.
func (me WithTimeout) TimeoutContext(ctx context.Context) (
	context.Context, context.CancelFunc,
) {
	return context.WithTimeout(ctx, me.Timeout)
}

// DB represents a database configuration
type DB struct {
	WithTimeout
	URI  string
	Name string
}

// CreateClient creates a client to connect to the DB.
func (me DB) CreateClient() (*mongo.Client, error) {
	return mongo.NewClient(
		options.Client().ApplyURI(me.URI),
	)
}

// Broker represents a brokwe configuration
type Broker struct {
	WithTimeout
	URI []string
}

// Connect to the server
func (me Broker) Connect() (*nats.Conn, error) {
	return nats.Connect(strings.Join(me.URI, ","), nats.Timeout(me.Timeout))
}

// Server represents a sever configuration.
type Server struct {
	WithTimeout
	Type string
	Addr string
}

// Config represents a configuration.
type Config struct {
	Db     DB
	Server Server
	Broker Broker
}
