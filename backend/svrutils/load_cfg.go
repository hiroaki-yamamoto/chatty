package svrutils

// CAUTION!! THIS MODULE CALLS fmt.Panicln THAT CAUSES TERMINATION ON ERROR

import (
	"log"
	"os"

	"github.com/hiroaki-yamamoto/real/backend/config"
)

// LoadCfg loads the configuration and returns the pointer of config.Config
func LoadCfg() *config.Config {
	cfg, err := config.New(os.Getenv("CONFIG_NAME"))
	if err != nil {
		log.Panicln("Loding Config Failed:", err)
		return nil
	}
	return cfg
}
