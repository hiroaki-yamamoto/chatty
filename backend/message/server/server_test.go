package server_test

import (
	"context"
	"fmt"
	"html"
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
		Expect(db.Collection(srvName).Drop(ctx)).Should(Succeed())
	})
	Describe("Subscription", func() {
		var models []*rpc.Message
		var ready sync.WaitGroup
		BeforeEach(func() {
			ready.Add(1)
			readyCh := make(chan *nats.Msg)
			sub, err := broker.ChanSubscribe("status/"+srvName+"/subscribe", readyCh)
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
			for count := 0; count < 45; count++ {
				countTxt := strconv.Itoa(count)
				msgToStream := &rpc.Message{
					SenderName: "<h1>Test Man</h1>" + countTxt,
					Message: fmt.Sprintf(
						"This is an %s example %d%s post from testman.",
						"<a href=\"https://example.com"+countTxt+"\">",
						count, "</a>",
					),
					Bump: count&0x01 == 0x01,
				}
				expMsg := msgToStream
				expMsg.SenderName = html.EscapeString(expMsg.SenderName)
				expMsg.Message = html.EscapeString(expMsg.Message)

				ready.Wait()
				status, err := cli.Post(subCli.Context(), &rpc.PostRequest{
					TopicId:   topicID.Hex(),
					Name:      msgToStream.SenderName,
					Message:   msgToStream.Message,
					Recaptcha: "PASSED",
					Bump:      msgToStream.Bump,
				})
				Expect(err).Should(Succeed())
				expMsg.Id = status.GetId()
				msg, err := subCli.Recv()
				expMsg.PostTime = msg.GetPostTime()
				Expect(err).Should(Succeed())
				Expect(msg).Should(Equal(expMsg))
			}
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
						SenderName: "<p>Test User " + numStr + "</p>",
						PostTime:   initPostDate.Add(time.Duration(i) * time.Hour),
						Message: `This is a <a href="javascript.alert('hello');">
            test</a>: ` + numStr,
						Host: "127.0.0.1",
						Bump: i&0x01 == 0x01,
					}
					cols[i] = model
					models[i] = &rpc.Message{
						Id:         model.ID.Hex(),
						SenderName: html.EscapeString(model.SenderName),
						PostTime: &timestamp.Timestamp{
							Seconds: model.PostTime.Unix(),
							Nanos:   int32((model.PostTime.Nanosecond() / 1000000) * 1000000),
						},
						Message: html.EscapeString(model.Message),
						Bump:    model.Bump,
					}
				}
				ctx, cancel := cfg.Db.TimeoutContext(context.Background())
				defer cancel()
				_, err := db.Collection(srvName).InsertMany(ctx, cols)
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
