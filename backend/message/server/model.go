package server

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model indicates the model of the message.
type Model struct {
	ID         primitive.ObjectID `bson:"_id"`
	TopicID    primitive.ObjectID
	SenderName string
	PostTime   time.Time
	Profile    string
	Message    string
	Host       string
}
