package main

import (
	"context"
	"log"
	"os"

	"github.com/hiroaki-yamamoto/real/backend/config"
)

func loadCfg() *config.Config {
	cfg, err := config.New(os.Getenv("CONFIG_NAME"))
	if err != nil {
		log.Panicln("Loding Config Failed:", err)
		return nil
	}
	return cfg
}

func connectDB(cfg *config.Config) {
	conCtx, cancelCon := cfg.Db.TimeoutContext(context.Background())
	defer cancelCon()
	if err := cfg.Db.Client.Connect(conCtx); err != nil {
		log.Panicln("Connecting to the DB Failed:", err)
	}
}

func constructSvr(cfg *config.Config) {

}

func main() {
	cfg := loadCfg()
	connectDB(cfg)
}
