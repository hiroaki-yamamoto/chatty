package server_test

import (
	"net"
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

const PKGNAME = "message"

var db *mongo.Database
var cli rpc.MessageServiceClient
var clicon *grpc.ClientConn
var cfg *config.Config
var lis net.Listener
var svr *grpc.Server

var _ = BeforeSuite(func() {
	cfg = svrutils.LoadCfg()
	cfg.Db.URI = "mongo://real:real@testdb"
	cfg.Server.Type = "unix"
	cfg.Server.Addr = "/tmp/" + PKGNAME + ".sock"
	db = svrutils.ConnectDB(cfg).Database(cfg.Db.Name)
	svr, lis = svrutils.Construct(cfg)

	go func() {
		if err := svr.Serve(lis); err != nil {
			Fail("Server Start Failed: " + err.Error())
		}
	}()
	if con, err := grpc.Dial(
		cfg.Server.Type+"://"+cfg.Server.Addr,
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
})
