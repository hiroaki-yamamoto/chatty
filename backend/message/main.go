package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

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

func disconnectDB(cfg *config.Config) {
	dbStopCtx, cancelDbStop := cfg.Db.TimeoutContext(context.Background())
	defer cancelDbStop()
	if err := cfg.Db.Client.Disconnect(dbStopCtx); err != nil {
		log.Panicln("Disconnecting the DB Failed:", err)
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

func trapInt(
	svr *grpc.Server,
	cfg *config.Config,
) (sig chan os.Signal) {
	sig = make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for range sig {
			log.Print("Gracefully Shutting The Server Down...")
			svr.GracefulStop()
			log.Print("Server Gracefuly Closed.")
		}
	}()
	return
}

func main() {
	cfg := loadCfg()
	connectDB(cfg)
	svr, lis := constructSvr(cfg)
	sig := trapInt(svr, cfg)
	defer close(sig)
	defer disconnectDB(cfg)
	if err := svr.Serve(lis); err != nil {
		log.Panicln("Server Start Failed: ", err)
	}
}
