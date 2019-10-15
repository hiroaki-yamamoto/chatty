package svrutils

import (
	"log"

	"github.com/hiroaki-yamamoto/real/backend/config"
	"github.com/nats-io/nats.go"
)

// CAUTION!! This package uses panic function

// InitBroker Initialized a broker.
func InitBroker(cfg *config.Config) (cli *nats.Conn) {
	var err error
	if cli, err = cfg.Broker.Connect(); err != nil {
		log.Panicln(err)
	}
	return
}
