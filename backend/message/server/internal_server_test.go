package server_test

import (
	"context"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	pr "go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/hiroaki-yamamoto/real/backend/message/server"
)

var _ = Describe("InternalServer", func() {
	Context("With Initial Model", func() {
		var topicIDs []pr.ObjectID
		BeforeEach(func() {
			const numTopics = 40
			const numResp = 40
			topicIDs = make([]pr.ObjectID, numTopics)
			msgs := make(bson.A, numTopics*numResp)
			for ti := 0; ti < numTopics; ti++ {
				topicID := pr.NewObjectID()
				topicIDs[ti] = topicID
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
				for ri := 0; ri < numResp; ri++ {
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
					msgs[(ri*numResp)+ri] = model
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
			topicIDs = nil
		})
		Context("Doesn't contain non-existent model", func() {
			It("Should recieve stats data", func() {
				statsCtx, cancelStats := context.WithTimeout(
					context.Background(), 10*time.Second,
				)
				defer cancelStats()
				// statsCli, err := prvCli.Stats(statsCtx)
				statsCli, err := prvCli.Stats(statsCtx)
				Expect(err).Should(Succeed())
				Fail("Not Implemented Yet")
			})
		})
	})
})
