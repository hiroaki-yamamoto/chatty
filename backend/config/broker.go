package config

import (
	"context"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// Broker represents a brokwe configuration
type Broker struct {
	Timeout time.Duration
	URI     []string
}

// TimeoutContext returns a context set a timeout with DB.Timeout.
func (me Broker) TimeoutContext(ctx context.Context) (
	context.Context, context.CancelFunc,
) {
	return context.WithTimeout(ctx, me.Timeout)
}

// Connect to the server
func (me Broker) Connect() (*nats.Conn, error) {
	return nats.Connect(strings.Join(me.URI, ","), nats.Timeout(me.Timeout))
}
