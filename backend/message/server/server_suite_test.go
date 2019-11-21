package server_test

import (
	"context"
	"net"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hiroaki-yamamoto/real/backend/config"
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

type SL struct {
	Server   *grpc.Server
	Listener net.Listener
}

const srvName = "messages"

var db *mongo.Database
var cli rpc.MessageServiceClient
var broker *nats.Conn
var clicon *grpc.ClientConn
var cfg *config.Config
var svrs []*SL

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = BeforeSuite(func() {
	setCfg()
	mockValidation()
	startBroker()
	db = svrutils.ConnectDB(cfg).Database(cfg.Db.Name)
	startSvr(cfg.Servers["message"])
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

func startSvr(svrCfg *config.Server) {
	svr, lis := svrutils.Construct(svrCfg)
	rpc.RegisterMessageServiceServer(
		svr, &server.Server{Setting: cfg, Database: db, Broker: broker},
	)
	svrs = append(svrs, &SL{
		Server:   svr,
		Listener: lis,
	})

	go func() {
		if err := svr.Serve(lis); err != nil {
			Fail("Server Start Failed: " + err.Error())
		}
	}()

	if con, err := grpc.Dial(svrCfg.Addr, grpc.WithInsecure()); err != nil {
		Fail("Connection Dial Failed: " + err.Error())
	} else {
		cli = rpc.NewMessageServiceClient(con)
		clicon = con
	}
}

var _ = AfterSuite(func() {
	if clicon != nil {
		clicon.Close()
	}
	for _, sl := range svrs {
		svr := sl.Server
		lis := sl.Listener
		if svr != nil {
			svr.GracefulStop()
		}
		if db != nil {
		}
		if lis != nil {
			lis.Close()
		}
	}
	svrutils.DisconnectDB(db.Client(), &cfg.Db)
	broker.Close()
})
