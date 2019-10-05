package svrutils

// CAUTION!! THIS MODULE CALLS fmt.Panicln THAT CAUSES TERMINATION ON ERROR

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/hiroaki-yamamoto/real/backend/config"
	"google.golang.org/grpc"
)

// Construct the server without calling Register**.
func Construct(cfg *config.Config) (*grpc.Server, net.Listener) {
	lis, err := net.Listen(cfg.Server.Type, cfg.Server.Addr)
	if err != nil {
		log.Panicln(err)
	}
	return grpc.NewServer(), lis
}

// Serve runs the server and trap int for graceful shutdown.
func Serve(
	lis net.Listener,
	svr *grpc.Server,
	cfg *config.Config,
) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for range sig {
			log.Print("Gracefully Shutting The Server Down...")
			svr.GracefulStop()
			log.Print("Server Gracefuly Closed.")
		}
	}()
	defer close(sig)
	if err := svr.Serve(lis); err != nil {
		log.Panicln("Server Start Failed: ", err)
	}
}
