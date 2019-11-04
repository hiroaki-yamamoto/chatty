package main

import (
	"github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/svrutils"
)

func main() {
	cfg := svrutils.LoadCfg()
	pubSvrCfg := cfg.Servers["message"]
	broker := svrutils.InitBroker(cfg)
	dbcli := svrutils.ConnectDB(cfg)
	svr, lis := svrutils.Construct(pubSvrCfg)
	rpc.RegisterMessageServiceServer(
		svr,
		&server.Server{
			Setting:  cfg,
			Database: dbcli.Database(cfg.Db.Name),
			Broker:   broker,
		},
	)
	svrutils.Serve(lis, svr, pubSvrCfg)
	defer svrutils.DisconnectDB(dbcli, &cfg.Db)
}
