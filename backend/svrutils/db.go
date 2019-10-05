package svrutils

// CAUTION!! THIS MODULE CALLS fmt.Panicln THAT CAUSES TERMINATION ON ERROR

import (
	"context"
	"log"

	"github.com/hiroaki-yamamoto/real/backend/config"
)

// ConnectDB attempts to connect to the DB.
func ConnectDB(cfg *config.Config) {
	conCtx, cancelCon := cfg.Db.TimeoutContext(context.Background())
	defer cancelCon()
	if err := cfg.Db.Client.Connect(conCtx); err != nil {
		log.Panicln("Connecting to the DB Failed:", err)
	}
}

// DisconnectDB attempts to disconnect to the DB.
func DisconnectDB(cfg *config.Config) {
	dbStopCtx, cancelDbStop := cfg.Db.TimeoutContext(context.Background())
	defer cancelDbStop()
	if err := cfg.Db.Client.Disconnect(dbStopCtx); err != nil {
		log.Panicln("Disconnecting the DB Failed:", err)
	}
}
