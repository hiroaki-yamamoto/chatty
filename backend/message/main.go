package main

import (
	"github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/svrutils"
)

func main() {
	cfg := svrutils.LoadCfg()
	svrutils.ConnectDB(cfg)
	svr, lis := svrutils.Construct(cfg)
	rpc.RegisterMessageServiceServer(svr, &server.Server{Setting: cfg})
	svrutils.Serve(lis, svr, cfg)
	defer svrutils.DisconnectDB(cfg)
}
