package server_test

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
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
	var models []*rpc.Message
	var topicID pr.ObjectID
	BeforeEach(func() {
		models = make([]*rpc.Message, 40)
		topicID = pr.NewObjectID()
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
					Nanos:   int32(model.PostTime.Nanosecond()),
				},
				Profile: model.Profile,
				Message: model.Message,
			}
		}
		ctx, cancel := cfg.Db.TimeoutContext(context.Background())
		defer cancel()
		_, err := db.Collection("messages").InsertMany(ctx, cols)
		Expect(err).Should(BeNil())
	})
	AfterEach(func() {
		ctx, cancel := cfg.Db.TimeoutContext(context.Background())
		models = nil
		topicID = pr.NilObjectID
		defer cancel()
		db.Collection("messages").Drop(ctx)
	})
	Describe("Subscription", func() {
		It("Reads the collection and return the docs initially", func() {
			ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
			var actual []*rpc.Message
			defer stop()
			subCli, err := cli.Subscribe(ctx, &rpc.MessageRequest{
				TopicId: topicID.Hex(),
			})
			Expect(err).Should(BeNil())
			for {
				msg, err := subCli.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).Should(BeNil())
				actual = append(actual, msg)
				if len(actual) >= len(models) {
					break
				}
			}
			Expect(actual).Should(Equal(models))
		})
	})
})
