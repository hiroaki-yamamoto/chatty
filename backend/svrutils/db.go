package svrutils

// CAUTION!! THIS MODULE CALLS fmt.Panicln THAT CAUSES TERMINATION ON ERROR

import (
	"context"
	"log"

	"github.com/hiroaki-yamamoto/real/backend/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// ConnectDB attempts to connect to the DB.
func ConnectDB(cfg *config.Config) (cli *mongo.Client) {
	conCtx, cancelCon := cfg.Db.TimeoutContext(context.Background())
	defer cancelCon()
	cli, err := cfg.Db.CreateClient()
	if err != nil {
		log.Panicln("Connecting to the DB Failed:", err)
	}
	if err = cli.Connect(conCtx); err != nil {
		log.Panicln("Connecting to the DB Failed:", err)
	}
	return
}

// DisconnectDB attempts to disconnect to the DB.
func DisconnectDB(cli *mongo.Client, cfg *config.DB) {
	dbStopCtx, cancelDbStop := cfg.TimeoutContext(context.Background())
	defer cancelDbStop()
	if err := cli.Disconnect(dbStopCtx); err != nil {
		log.Panicln("Disconnecting the DB Failed:", err)
	}
}
