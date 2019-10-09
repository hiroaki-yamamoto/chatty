package server_test

import (
	"log"
	"testing"

	"github.com/hiroaki-yamamoto/real/backend/config"
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

var cli *mongo.Client
var cfg *config.Config
var svr *grpc.Server

var _ = BeforeSuite(func() {
	cfg = svrutils.LoadCfg()
	cfg.Db.URI = "mongo://real:real@testdb"
	cfg.Server.Type = "unix"
	cfg.Server.Addr = "/tmp/real-test.sock"
	cli = svrutils.ConnectDB(cfg)
	svr, lis := svrutils.Construct(cfg)

	go func() {
		if err := svr.Serve(lis); err != nil {
			log.Panicln("Server Start Failed: ", err)
		}
	}()
})

var _ = AfterSuite(func() {
	svr.GracefulStop()
	svrutils.DisconnectDB(cli, &cfg.Db)
})
