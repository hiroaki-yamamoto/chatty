package server_test

import (
	"context"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	pr "go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/random"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
)

const randomCharMap = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var _ = Describe("Message Server", func() {
	var topicID pr.ObjectID
	BeforeEach(func() {
		topicID = pr.NewObjectID()
	})
	AfterEach(func() {
		ctx, cancel := cfg.Db.TimeoutContext(context.Background())
		defer cancel()
		Expect(db.Collection("messages").Drop(ctx)).Should(Succeed())
	})
	Describe("Subscription", func() {
		var models []*rpc.Message
		var ready sync.WaitGroup
		BeforeEach(func() {
			ready.Add(1)
			readyCh := make(chan *nats.Msg)
			sub, err := broker.ChanSubscribe("status/messages/subscribe", readyCh)
			Expect(err).Should(Succeed())
			go func() {
				defer close(readyCh)
				status := <-readyCh
				switch string(status.Data) {
				case "ready":
					sub.Unsubscribe()
					ready.Done()
				}
			}()
		})
		checkPostMsg := func(subCli rpc.MessageService_SubscribeClient) {
			additionalPostTime := time.Now().UTC().Add(-240 * time.Hour)
			msgToStream := &rpc.Message{
				SenderName: "Test Man",
				PostTime: &timestamp.Timestamp{
					Seconds: additionalPostTime.Unix(),
					Nanos: int32(
						(additionalPostTime.Nanosecond() / 1000000) * 1000000,
					),
				},
				Profile: "https://google.com",
				Message: "This is an example post from testman.",
			}

			ready.Wait()
			status, err := cli.Post(subCli.Context(), &rpc.PostRequest{
				TopicId:   topicID.Hex(),
				Name:      msgToStream.SenderName,
				Message:   msgToStream.Message,
				Recaptcha: "PASSED",
			})
			Expect(err).Should(Succeed())
			msgToStream.Id = status.Id
			msg, err := subCli.Recv()
			Expect(err).Should(Succeed())
			Expect(msg).Should(Equal(msgToStream))
		}
		checkInitMsg := func() (
			rpc.MessageService_SubscribeClient,
			context.CancelFunc,
		) {
			ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
			actual := make([]*rpc.Message, cap(models))
			subCli, err := cli.Subscribe(ctx, &rpc.MessageRequest{
				TopicId: topicID.Hex(),
			})
			Expect(err).Should(Succeed())
			for i := 0; i < cap(actual); i++ {
				msg, err := subCli.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).Should(Succeed())
				actual[i] = msg
			}
			Expect(actual).Should(Equal(models))
			return subCli, stop
		}
		Context("With initial messages", func() {
			BeforeEach(func() {
				models = make([]*rpc.Message, 48)
				cols := make(bson.A, cap(models))
				initPostDate := time.Now().UTC().Add(
					-time.Duration(cap(models)) * time.Hour,
				)
				for i := 0; i < cap(models); i++ {
					numStr := strconv.Itoa(i)
					msgID := pr.NewObjectID()
					model := &Model{
						ID:         msgID,
						TopicID:    topicID,
						SenderName: "Test User " + numStr,
						PostTime:   initPostDate.Add(time.Duration(i) * time.Hour),
						Profile:    random.GenerateRandomText(randomCharMap, 16),
						Message:    "This is a test: " + numStr,
						Host:       "127.0.0.1",
					}
					cols[i] = model
					models[i] = &rpc.Message{
						Id:         model.ID.Hex(),
						SenderName: model.SenderName,
						PostTime: &timestamp.Timestamp{
							Seconds: model.PostTime.Unix(),
							Nanos:   int32((model.PostTime.Nanosecond() / 1000000) * 1000000),
						},
						Profile: model.Profile,
						Message: model.Message,
					}
				}
				ctx, cancel := cfg.Db.TimeoutContext(context.Background())
				defer cancel()
				_, err := db.Collection("messages").InsertMany(ctx, cols)
				Expect(err).Should(Succeed())
			})
			AfterEach(func() {
				models = nil
			})
			It("Reads the collection and return the docs initially", func() {
				_, stop := checkInitMsg()
				defer stop()
			})
			It("Recives the message when it's posted.", func() {
				subCli, stop := checkInitMsg()
				defer stop()
				checkPostMsg(subCli)
			}, 8)
		})
		Context("Without any initial messages", func() {
			It("Recives the message when it's posted.", func() {
				ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
				defer stop()
				subCli, err := cli.Subscribe(ctx, &rpc.MessageRequest{
					TopicId: topicID.Hex(),
				})
				Expect(err).Should(Succeed())

				checkPostMsg(subCli)
			}, 8)
		})
	})
})
