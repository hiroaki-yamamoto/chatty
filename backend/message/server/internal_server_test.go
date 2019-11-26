package server_test

import (
	"context"
	"strconv"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	pr "go.mongodb.org/mongo-driver/bson/primitive"

	intRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	. "github.com/hiroaki-yamamoto/real/backend/message/server"
)

var _ = Describe("InternalServer", func() {
	Context("With Initial Model", func() {
		var expMsgs []*intRPC.StatsResponse
		BeforeEach(func() {
			const numTopics = 40
			const numResp = 40
			expMsgs = make([]*intRPC.StatsResponse, numTopics)
			var msgs bson.A
			for ti := 0; ti < numTopics; ti++ {
				topicID := pr.NewObjectID()
				expMsgs[ti] = &intRPC.StatsResponse{
					TopicId: topicID.Hex(),
					NumMsgs: int64(numResp - ti),
				}
				initPostDate := time.Now().UTC()
				initPostDate = initPostDate.Add(
					-time.Duration(initPostDate.Minute()) * time.Minute,
				).Add(
					-time.Duration(initPostDate.Second()) * time.Second,
				).Add(
					-time.Duration(initPostDate.Nanosecond()) * time.Nanosecond,
				).Add(
					time.Duration(ti) * time.Minute,
				)
				for ri := ti; ri < numResp; ri++ {
					msgNum := strconv.Itoa(ri)
					if len(msgNum) < 1 {
						msgNum = " **( Undefined )** "
					}
					model := &Model{
						ID:         pr.NewObjectID(),
						TopicID:    topicID,
						SenderName: "<p>Test User " + msgNum + "</p>",
						PostTime:   initPostDate.Add(time.Duration(ri) * time.Second),
						Message: `This is a <a href="javascript.alert('hello');">
            test</a>: ` + msgNum,
						Host: "127.0.0.1",
						Bump: ri&0x01 == 0x01,
					}
					msgs = append(msgs, model)
				}
			}
			ctx, cancel := cfg.Db.TimeoutContext(context.Background())
			defer cancel()
			_, err := db.Collection(srvName).InsertMany(ctx, msgs)
			Expect(err).Should(Succeed())
		})
		AfterEach(func() {
			ctx, cancel := cfg.Db.TimeoutContext(context.Background())
			defer cancel()
			err := db.Collection(srvName).Drop(ctx)
			Expect(err).Should(Succeed())
			expMsgs = nil
		})
		Context("Doesn't contain non-existent model", func() {
			It("Should recieve stats data", func() {
				statsCtx, cancelStats := context.WithTimeout(
					context.Background(), 10*time.Second,
				)
				defer cancelStats()
				statsCli, err := prvCli.Stats(statsCtx)
				Expect(err).Should(Succeed())

				statsLst := make([]*intRPC.StatsResponse, len(expMsgs))
				var wg sync.WaitGroup
				sendErrCh := make(chan error)
				recvErrCh := make(chan error)
				wg.Add(2)
				go func(errCh chan error) {
					for _, expMsg := range expMsgs {
						errCh <- statsCli.Send(
							&intRPC.StatsRequest{TopicId: expMsg.TopicId},
						)
					}
					wg.Done()
				}(sendErrCh)
				go func(errCh chan error) {
					for i := 0; i < len(statsLst); i++ {
						stats, err := statsCli.Recv()
						errCh <- err
						if err != nil {
							return
						}
						statsLst[i] = stats
					}
					wg.Done()
				}(recvErrCh)
				for i := 0; i < len(expMsgs); i++ {
					Expect(<-sendErrCh).Should(Succeed())
					Expect(<-recvErrCh).Should(Succeed())
				}
				wg.Wait()
				Expect(len(statsLst)).Should(Equal(len(expMsgs)))
				Expect(statsLst).Should(ConsistOf(expMsgs))
			}, 4000)
		})
	})
})
