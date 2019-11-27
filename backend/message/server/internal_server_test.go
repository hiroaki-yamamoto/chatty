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
		var statsCli intRPC.MessageStats_StatsClient
		var stopStats context.CancelFunc
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
			ctx, cancelInsert := cfg.Db.TimeoutContext(context.Background())
			defer cancelInsert()
			_, err := db.Collection(srvName).InsertMany(ctx, msgs)
			Expect(err).Should(Succeed())

			var statsCtx context.Context
			statsCtx, stopStats = context.WithTimeout(
				context.Background(), 10*time.Second,
			)
			statsCli, err = prvCli.Stats(statsCtx)
			Expect(err).Should(Succeed())
		})
		AfterEach(func() {
			stopStats()
			ctx, cancel := cfg.Db.TimeoutContext(context.Background())
			defer cancel()
			err := db.Collection(srvName).Drop(ctx)
			Expect(err).Should(Succeed())
			expMsgs = nil
		})
		Context("Doesn't contain non-existent model", func() {
			collect := func() (resp []*intRPC.StatsResponse, err error) {
				var wg sync.WaitGroup
				wg.Add(2)
				sendErrCh := make(chan error)
				recvErrCh := make(chan error)
				go func() {
					defer wg.Done()
					defer close(sendErrCh)
					for _, expMsg := range expMsgs {
						err := statsCli.Send(
							&intRPC.StatsRequest{TopicId: expMsg.TopicId},
						)
						if err != nil {
							sendErrCh <- err
							return
						}
					}
				}()
				go func() {
					defer wg.Done()
					defer close(recvErrCh)
					for i := 0; i < len(expMsgs); i++ {
						stats, err := statsCli.Recv()
						if err != nil {
							recvErrCh <- err
							return
						}
						resp = append(resp, stats)
					}
				}()
				if sendErr, opened := <-sendErrCh; sendErr != nil && opened {
					err = sendErr
				}
				if recvErr, opened := <-recvErrCh; recvErr != nil && opened {
					err = recvErr
				}
				wg.Wait()
				return
			}
			It("Should recieve stats data", func() {
				statsLst, err := collect()
				Expect(err).Should(Succeed())
				Expect(len(statsLst)).Should(Equal(len(expMsgs)))
				Expect(statsLst).Should(ConsistOf(expMsgs))
			})
			// It("Should receive continuously", func() {
			// })
		})
	})
})
