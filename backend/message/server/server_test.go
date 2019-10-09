package server_test

import (
	"context"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"

	. "github.com/hiroaki-yamamoto/real/backend/message/server"
	"github.com/hiroaki-yamamoto/real/backend/random"
)

const randomCharMap = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var _ = Describe("Message Server", func() {
	var models []Model
	BeforeEach(func() {
		models = make([]Model, 40)
		ctx, cancel := cfg.Db.TimeoutContext(context.Background())
		defer cancel()
		cols := make(bson.A, cap(models))
		for i := 0; i < cap(models); i++ {
			numStr := strconv.Itoa(i)
			models[i] = Model{
				SenderName: "Test User " + numStr,
				PostTime:   time.Now().UTC().Add(-time.Duration(i) * time.Hour),
				Profile:    random.GenerateRandomText(randomCharMap, 16),
				Message:    "This is a test: " + numStr,
				Host:       "127.0.0.1",
			}
			cols[i] = &models[i]
		}
		db.Collection("messages").InsertMany(ctx, cols)
	})
	AfterEach(func() {
		ctx, cancel := cfg.Db.TimeoutContext(context.Background())
		models = nil
		defer cancel()
		db.Collection("messages").Drop(ctx)
	})
	Describe("Subscription", func() {
		It("Reads the collection and return the docs initially", func() {
			// How can I write this???
		})
	})
})
