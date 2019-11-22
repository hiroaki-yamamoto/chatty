package main

import (
	statsRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	"github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/svrutils"
)

func main() {
	manager := svrutils.ServerManager{}
	cfg := svrutils.LoadCfg()
	broker := svrutils.InitBroker(cfg)
	dbcli := svrutils.ConnectDB(cfg)
	defer broker.Close()
	defer manager.CloseAll()
	defer svrutils.DisconnectDB(dbcli, &cfg.Db)

	svr, _ := manager.Construct(cfg.Servers["message"])
	statsSvr, _ := manager.Construct(cfg.Servers["message/stats"])
	rpc.RegisterMessageServiceServer(
		svr,
		&server.Server{
			Setting:  cfg,
			Database: dbcli.Database(cfg.Db.Name),
			Broker:   broker,
		},
	)
	statsRPC.RegisterMessageStatsServer(statsSvr, &server.InternalServer{
		DB:     dbcli.Database(cfg.Db.Name),
		Broker: broker,
	})
	go manager.TrapInt()
	manager.Serve()
}
