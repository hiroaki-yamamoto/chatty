package svrutils

// CAUTION!! THIS MODULE CALLS fmt.Panicln THAT CAUSES TERMINATION ON ERROR

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/hiroaki-yamamoto/real/backend/config"
	"google.golang.org/grpc"
)

type server struct {
	Server   *grpc.Server
	Listener net.Listener
}

// ServerManager is used for server management.
type ServerManager struct {
	Servers []*server
}

// Construct the server and listener.
func (me *ServerManager) Construct(
	cfg *config.Server,
	opts ...grpc.ServerOption,
) (*grpc.Server, net.Listener) {
	lis, err := net.Listen(cfg.Type, cfg.Addr)
	if err != nil {
		log.Panicln(err)
	}
	svr := &server{
		Server:   grpc.NewServer(opts...),
		Listener: lis,
	}
	me.Servers = append(me.Servers, svr)
	return svr.Server, svr.Listener
}

// TrapInt handles Graceful Stop by SIGINT.
// Note that this function runs **synchronously**, not async.
func (me *ServerManager) TrapInt() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	defer close(sig)
	for range sig {
		log.Print("Gracefully Shutting The Server Down...")
		me.CloseAll()
		log.Print("Server Gracefuly Closed.")
	}
}

// CloseAll closes all the servers and listeners.
func (me *ServerManager) CloseAll() {
	for _, srv := range me.Servers {
		srv.Server.GracefulStop()
		srv.Listener.Close()
	}
	me.Servers = nil
}

// Serve runs the server.
func (me *ServerManager) Serve() {
	var wg sync.WaitGroup
	wg.Add(len(me.Servers))
	for _, srv := range me.Servers {
		addr := srv.Listener.Addr()
		log.Printf(
			"Opening the server on %s as %s socket\n",
			addr.String(), addr.Network(),
		)
		go func(svr *grpc.Server, lis net.Listener) {
			if err := svr.Serve(lis); err != nil {
				log.Panicln("Server Start Failed: ", err)
			}
			wg.Done()
		}(srv.Server, srv.Listener)
	}
	wg.Wait()
}
