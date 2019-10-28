package server

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hiroaki-yamamoto/real/backend/config"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	topicID, err := primitive.ObjectIDFromHex(req.TopicId)
	if err != nil {
		return
	}

	col := me.Database.Collection("messages")
	query := bson.M{"topicid": topicID}
	findCur, err := col.Find(
		stream.Context(), query,
		&options.FindOptions{
			Skip: &start,
			Sort: bson.M{
				"posttime": 1,
			},
		},
	)
	if err != nil {
		return
	}

	for findCur.Next(stream.Context()) {
		var model Model
		if err = findCur.Decode(&model); err != nil {
			return
		}
		err = stream.Send(&rpc.Message{
			Id:         model.ID.Hex(),
			SenderName: model.SenderName,
			PostTime: &timestamp.Timestamp{
				Seconds: model.PostTime.Unix(),
				Nanos:   int32(model.PostTime.Nanosecond()),
			},
			Message: model.Message,
			Profile: model.Profile,
		})
		if err != nil {
			return
		}
	}

	msgCh := make(chan *nats.Msg)
	defer close(msgCh)
	chSub, err := me.Broker.ChanSubscribe("messages/"+req.TopicId, msgCh)
	if err != nil {
		return
	}
	defer chSub.Unsubscribe()
	me.Broker.Publish("ready", nil)
	for {
		select {
		case msg := <-msgCh:
			var model rpc.Message
			if err = msgpack.Unmarshal(msg.Data, &model); err != nil {
				return
			}
			stream.Send(&model)
			break
		case <-stream.Context().Done():
			return
		}
	}
}

// Post recirds the message, broadcast it, and returns Status structure
func (me *Server) Post(req *rpc.PostRequest) (*rpc.Status, error) {

}
