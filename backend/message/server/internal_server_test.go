package server_test

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	pr "go.mongodb.org/mongo-driver/bson/primitive"

	intRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	. "github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
)

var _ = Describe("InternalServer", func() {
	Context("With Initial Model", func() {
		const numTopics = 40
		const numResp = 40
		var expMsgs []*intRPC.StatsResponse
		var statsCli intRPC.MessageStats_StatsClient
		var stopStats context.CancelFunc
		BeforeEach(func() {
			expMsgs = make([]*intRPC.StatsResponse, numTopics)
			var msgs bson.A
			for ti := 0; ti < numTopics; ti++ {
				topicID := pr.NewObjectID()
				initPostDate := time.Now().UTC()
				initPostDate = initPostDate.Add(
					(-1 * time.Hour) -
						(time.Duration(initPostDate.Hour()) * time.Hour) -
						(time.Duration(initPostDate.Minute()) * time.Minute) -
						(time.Duration(initPostDate.Second()) * time.Second),
				).Add(
					time.Duration(ti) * time.Minute,
				)
				lastBump := initPostDate
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
					if model.Bump && model.PostTime.After(lastBump) {
						lastBump = model.PostTime
					}
				}
				expMsgs[ti] = &intRPC.StatsResponse{
					TopicId: topicID.Hex(),
					NumMsgs: int64(numResp - ti),
					LastBump: &timestamp.Timestamp{
						Seconds: lastBump.Unix(),
						Nanos:   int32((lastBump.Nanosecond() / 1000000) * 1000000),
					},
				}
			}
			ctx, cancelInsert := cfg.Db.TimeoutContext(context.Background())
			defer cancelInsert()
			_, err := db.Collection(srvName).InsertMany(ctx, msgs)
			Expect(err).Should(Succeed())

			var statsCtx context.Context
			statsCtx, stopStats = context.WithTimeout(
				context.Background(), 3*time.Second,
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
		Describe("Request all the model", func() {
			var statsLst []*intRPC.StatsResponse
			BeforeEach(func() {
				var err error
				sendErrCh := make(chan error)
				recvErrCh := make(chan error)
				go func() {
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
					defer close(recvErrCh)
					for i := 0; i < len(expMsgs); i++ {
						var stats *intRPC.StatsResponse
						stats, err = statsCli.Recv()
						if err != nil {
							recvErrCh <- err
						}
						statsLst = append(statsLst, stats)
					}
					return
				}()
				if recvErr, opened := <-recvErrCh; opened {
					Expect(recvErr).Should(Succeed())
				}
				if sendErr, opened := <-sendErrCh; opened {
					Expect(sendErr).Should(Succeed())
				}
			})
			Context("Without non-existence topic ID request", func() {
				It("Should recieve stats data", func() {
					num := len(expMsgs)
					Expect(len(statsLst)).Should(Equal(num))
					Expect(statsLst).Should(ConsistOf(expMsgs))
				})
				It("Should receive continuously", func() {
					targetMsg := expMsgs[rand.Intn(len(expMsgs))]
					sendErrCh := make(chan error)
					recvErrCh := make(chan error)
					statsCh := make(chan *intRPC.StatsResponse)

					go func() {
						defer close(sendErrCh)
						ctx, cancel := context.WithTimeout(
							context.Background(), 3*time.Second,
						)
						defer cancel()
						_, err := pubCli.Post(ctx, &rpc.PostRequest{
							TopicId:   targetMsg.GetTopicId(),
							Name:      "<p>Test User </p>",
							Message:   `This is a <a href="javascript.alert('hello');">test</a>`,
							Recaptcha: "PASSED",
							Bump:      true,
						})
						sendErrCh <- err
					}()
					go func() {
						defer close(recvErrCh)
						defer close(statsCh)
						stats, err := statsCli.Recv()
						recvErrCh <- err
						statsCh <- stats
					}()
					Expect(<-sendErrCh).Should(Succeed())
					Expect(<-recvErrCh).Should(Succeed())
					targetMsg.NumMsgs++
					stats := <-statsCh
					Expect(time.Unix(
						stats.GetLastBump().GetSeconds(),
						int64(stats.GetLastBump().GetNanos()),
					)).Should(BeTemporally(
						">", time.Unix(
							targetMsg.GetLastBump().GetSeconds(),
							int64(targetMsg.GetLastBump().GetNanos()),
						),
					))
					stats.LastBump = targetMsg.GetLastBump()
					Expect(stats).Should(Equal(targetMsg))
				})
			})
			Context("With non-existence topic ID request", func() {
				nonExistence := make([]pr.ObjectID, numTopics)
				BeforeEach(func() {
					for i := 0; i < len(nonExistence); i++ {
						topicID := pr.NewObjectID()
						err := statsCli.Send(
							&intRPC.StatsRequest{TopicId: topicID.Hex()},
						)
						Expect(err).Should(Succeed())
						nonExistence[i] = topicID
					}
				})
				It("Shouldn't receive any extra fields", func() {
					recvCh := make(chan *intRPC.StatsResponse)
					errCh := make(chan error)
					defer close(recvCh)
					defer close(errCh)
					go func() {
						stats, err := statsCli.Recv()
						if _, open := <-errCh; err != nil && open {
							errCh <- err
						}
						if _, open := <-recvCh; stats != nil && open {
							recvCh <- stats
						}
					}()
					for {
						select {
						case stat := <-recvCh:
							Expect(stat).To(BeNil())
						case err := <-errCh:
							Expect(err).To(Succeed())
						case <-time.After(100 * time.Millisecond):
							return
						}
					}
				})
			})
		})
	})
})
