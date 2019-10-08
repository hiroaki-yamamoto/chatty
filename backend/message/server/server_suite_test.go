package server_test

import (
	"testing"

	"github.com/hiroaki-yamamoto/real/backend/svrutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = BeforeSuite(func() {
	cfg := svrutils.LoadCfg()
})
