package svrutils

import (
	"log"
	"os"
	"os/signal"

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
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)
		defer close(sig)
		for range sig {
			log.Print("Closing Broker Client...")
			cli.Close()
			log.Print("Broker Client has been closed...")
		}
	}()
	return
}
