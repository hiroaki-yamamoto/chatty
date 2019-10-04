package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/hiroaki-yamamoto/real/backend/config"
	"github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"google.golang.org/grpc"
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

func constructSvr(cfg *config.Config) (*grpc.Server, net.Listener) {
	lis, err := net.Listen(cfg.Server.Type, cfg.Server.Addr)
	if err != nil {
		log.Panicln(err)
	}
	svr := grpc.NewServer()
	rpc.RegisterMessageServiceServer(svr, &server.Server{Setting: cfg})
	return svr, lis
}

func main() {
	cfg := loadCfg()
	connectDB(cfg)
	svr, lis := constructSvr(cfg)
	if err := svr.Serve(lis); err != nil {
		log.Panicln(err)
	}
}
