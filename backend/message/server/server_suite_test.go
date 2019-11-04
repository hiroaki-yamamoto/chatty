package server_test

import (
	"net"
	"testing"

	"github.com/hiroaki-yamamoto/real/backend/config"
	"github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/svrutils"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var addr = "localhost:50000"
var db *mongo.Database
var cli rpc.MessageServiceClient
var broker *nats.Conn
var clicon *grpc.ClientConn
var cfg *config.Config
var lis net.Listener
var svr *grpc.Server

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = BeforeSuite(func() {
	cfg = svrutils.LoadCfg()
	cfg.Db.URI = "mongodb://real:real@testdb/"
	cfg.Broker.URI = []string{"nats://testbroker:4222"}
	pubSvrCfg := cfg.Servers["message"]
	pubSvrCfg.Addr = addr
	db = svrutils.ConnectDB(cfg).Database(cfg.Db.Name)
	svr, lis = svrutils.Construct(pubSvrCfg)
	broker = svrutils.InitBroker(cfg)
	rpc.RegisterMessageServiceServer(
		svr, &server.Server{Setting: cfg, Database: db, Broker: broker},
	)

	go func() {
		if err := svr.Serve(lis); err != nil {
			Fail("Server Start Failed: " + err.Error())
		}
	}()
	if con, err := grpc.Dial(
		pubSvrCfg.Addr,
		grpc.WithInsecure(),
	); err != nil {
		Fail("Connection Dial Failed: " + err.Error())
	} else {
		cli = rpc.NewMessageServiceClient(con)
		clicon = con
	}
})

var _ = AfterSuite(func() {
	if clicon != nil {
		clicon.Close()
	}
	if svr != nil {
		svr.GracefulStop()
	}
	if db != nil {
	}
	if lis != nil {
		lis.Close()
	}
	svrutils.DisconnectDB(db.Client(), &cfg.Db)
	broker.Close()
})
