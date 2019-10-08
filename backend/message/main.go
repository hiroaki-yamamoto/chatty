package main

import (
	"github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/svrutils"
)

func main() {
	cfg := svrutils.LoadCfg()
	dbcli := svrutils.ConnectDB(cfg)
	svr, lis := svrutils.Construct(cfg)
	rpc.RegisterMessageServiceServer(
		svr, &server.Server{Setting: cfg, Database: dbcli.Database(cfg.Db.Name)},
	)
	svrutils.Serve(lis, svr, cfg)
	defer svrutils.DisconnectDB(dbcli, &cfg.Db)
}
