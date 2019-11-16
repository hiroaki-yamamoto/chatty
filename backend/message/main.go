package main

import (
	statsRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	"github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/svrutils"
	"sync"
)

func main() {
	var svrLock sync.WaitGroup
	svrLock.Add(2) // main.go will launch 2 servers; public RPC and stats.
	cfg := svrutils.LoadCfg()
	pubSvrCfg := cfg.Servers["message"]
	statsSvrCfg := cfg.Servers["message/stats"]
	broker := svrutils.InitBroker(cfg)
	dbcli := svrutils.ConnectDB(cfg)
	defer svrutils.DisconnectDB(dbcli, &cfg.Db)
	svr, lis := svrutils.Construct(pubSvrCfg)
	statsSvr, statsLis := svrutils.Construct(statsSvrCfg)
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
	go func() {
		svrutils.Serve(lis, svr, pubSvrCfg)
		svrLock.Done()
	}()
	go func() {
		svrutils.Serve(statsLis, statsSvr, statsSvrCfg)
		svrLock.Done()
	}()
	svrLock.Wait()
}
