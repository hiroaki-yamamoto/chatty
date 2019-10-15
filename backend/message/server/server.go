package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hiroaki-yamamoto/real/backend/config"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Server implements MessageServiceServer interface.
type Server struct {
	Setting  *config.Config
	Database *mongo.Database
	Broker   *nats.Conn
}

// Subscribe handles subscribtions from users
func (me *Server) Subscribe(
	req *rpc.MessageRequest, stream rpc.MessageService_SubscribeServer,
) (err error) {
	start := int64(req.StartFrom)
	col := me.Database.Collection("messages")
	findCtx, cancelFind := me.Setting.Db.TimeoutContext(stream.Context())
	defer cancelFind()
	query := bson.M{"topicId": req.TopicId}
	findCur, err := col.Find(
		findCtx, query, &options.FindOptions{Skip: &start},
	)
	if err != nil {
		return
	}
	chstream, err := col.Watch(
		stream.Context(),
		bson.A{bson.M{"$match": query}},
	)
	if err != nil {
		return
	}

	decode := func(cur interface {
		Next(context.Context) bool
		Decode(interface{}) error
	}) {
		for nxtCtx, stopNxt := me.Setting.Db.TimeoutContext(
			stream.Context(),
		); cur.Next(nxtCtx); stopNxt() {
			var model Model
			if err = cur.Decode(&model); err != nil {
				return
			}
			err = stream.Send(&rpc.Message{
				Id:         model.ID.String(),
				SenderName: model.SenderName,
				PostTime: &timestamp.Timestamp{
					Seconds: model.PostTime.Unix(),
					Nanos:   int32(model.PostTime.Nanosecond()),
				},
			})
			if err != nil {
				return
			}
		}
	}
	decode(findCur)
	decode(chstream)
	defer chstream.Close(stream.Context())
	return
}
