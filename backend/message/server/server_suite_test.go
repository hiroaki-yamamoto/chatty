package server_test

import (
	"log"
	"testing"

	"github.com/hiroaki-yamamoto/real/backend/config"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/svrutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var db *mongo.Client
var cli rpc.MessageServiceClient
var clicon *grpc.ClientConn
var cfg *config.Config
var svr *grpc.Server

var _ = BeforeSuite(func() {
	cfg = svrutils.LoadCfg()
	cfg.Db.URI = "mongo://real:real@testdb"
	cfg.Server.Type = "unix"
	cfg.Server.Addr = "/tmp/real-test.sock"
	db = svrutils.ConnectDB(cfg)
	svr, lis := svrutils.Construct(cfg)

	go func() {
		if err := svr.Serve(lis); err != nil {
			log.Panicln("Server Start Failed: ", err)
		}
	}()
	if con, err := grpc.Dial(cfg.Server.Addr); err != nil {
		log.Panicln("Connection Dial Failed: ", err)
	} else {
		cli = rpc.NewMessageServiceClient(con)
		clicon = con
	}
})

var _ = AfterSuite(func() {
	clicon.Close()
	svr.GracefulStop()
	svrutils.DisconnectDB(db, &cfg.Db)
})
