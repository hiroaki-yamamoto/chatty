package server_test

import (
	"context"
	"net"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hiroaki-yamamoto/real/backend/config"
	intRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	"github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/svrutils"
	"github.com/hiroaki-yamamoto/real/backend/validation"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

const srvName = "messages"

var db *mongo.Database
var pubCli rpc.MessageServiceClient
var prvCli intRPC.MessageStatsClient
var broker *nats.Conn
var cfg *config.Config
var svrMgr svrutils.ServerManager
var clicons []*grpc.ClientConn

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = BeforeSuite(func() {
	setCfg()
	mockValidation()
	startBroker()
	db = svrutils.ConnectDB(cfg).Database(cfg.Db.Name)
	pubSvr, pubLis := svrMgr.Construct(cfg.Servers["message"])
	intSvr, intLis := svrMgr.Construct(cfg.Servers["message/stats"])
	preparePubServer(pubSvr, pubLis)
	prepareInternalServer(intSvr, intLis)
	go svrMgr.TrapInt()
	go svrMgr.Serve()
})

func setCfg() {
	cfg = svrutils.LoadCfg()
	cfg.Db.URI = "mongodb://real:real@testdb/"
	cfg.Broker.URI = []string{"nats://testbroker:4222"}
	cfg.Servers["message"].Addr = "localhost:50000"
	cfg.Servers["message/stats"].Addr = "localhost:50001"
}

func mockValidation() {
	validation.New = func(
		reqCtx context.Context,
		recapSecret string,
	) (*validator.Validate, error) {
		vld := validator.New()
		vld.RegisterValidation("recap", func(fl validator.FieldLevel) bool {
			return true
		})
		return vld, nil
	}
}

func startBroker() {
	broker = svrutils.InitBroker(cfg)
}

func preparePubServer(svr *grpc.Server, lis net.Listener) {
	addr := lis.Addr()
	rpc.RegisterMessageServiceServer(
		svr, &server.Server{Setting: cfg, Database: db, Broker: broker},
	)

	if con, err := grpc.Dial(addr.String(), grpc.WithInsecure()); err != nil {
		Fail("Connection Dial Failed: " + err.Error())
	} else {
		pubCli = rpc.NewMessageServiceClient(con)
		clicons = append(clicons, con)
	}
}

func prepareInternalServer(svr *grpc.Server, lis net.Listener) {
	addr := lis.Addr()
	intRPC.RegisterMessageStatsServer(
		svr, &server.InternalServer{DB: db, Broker: broker},
	)

	if con, err := grpc.Dial(addr.String(), grpc.WithInsecure()); err != nil {
		Fail("Connection Dial Failed: " + err.Error())
	} else {
		prvCli = intRPC.NewMessageStatsClient(con)
		clicons = append(clicons, con)
	}
}

var _ = AfterSuite(func() {
	svrMgr.CloseAll()
	for _, con := range clicons {
		if con != nil {
			con.Close()
		}
	}
	clicons = nil
	svrutils.DisconnectDB(db.Client(), &cfg.Db)
	broker.Close()
})
