package server

import (
	"context"
	"html"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Model indicates the model of the message.
type Model struct {
	ID         primitive.ObjectID `bson:"_id"`
	TopicID    primitive.ObjectID `validate:"required"`
	SenderName string
	PostTime   time.Time `validate:"required"`
	Message    string    `validate:"required"`
	Host       string    `validate:"required"`
	Recaptcha  string    `bson:"-" validate:"recap"`
}

// Store the model to the collection.
// WARNING: this doesn't update the model that already exists. The behavior
//  is only inserting the model.
func (me *Model) Store(
	ctx context.Context,
	col *mongo.Collection,
) (err error) {
	var res *mongo.InsertOneResult
	res, err = col.InsertOne(ctx, me)
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("[WARN] The returned ID was not ObjectID.")
	}
	me.ID = id
	return
}

// ToRPCMsg converts the model into *rpc.Message
func (me *Model) ToRPCMsg(escape bool) *rpc.Message {
	ret := &rpc.Message{
		Id:         me.ID.Hex(),
		SenderName: me.SenderName,
		PostTime: &timestamp.Timestamp{
			Seconds: me.PostTime.Unix(),
			Nanos:   int32(me.PostTime.Nanosecond()),
		},
		Message: me.Message,
	}
	if escape {
		ret.SenderName = html.EscapeString(me.SenderName)
		ret.Message = html.EscapeString(me.Message)
	}
	return ret
}
