package message

import (
	"github.com/hiroaki-yamamoto/real/backend/config"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Server implements MessageServiceServer interface.
type Server struct {
	cfg config.Config
}

// Subscribe handles subscribtions from users
func (me *Server) Subscribe(
	req *rpc.MessageRequest, stream rpc.MessageService_SubscribeServer,
) (err error) {
	start := int64(req.StartFrom)
	col := me.cfg.Db.Database.Collection("messages")
	findCtx, cancelFind := me.cfg.Db.TimeoutContext(stream.Context())
	defer cancelFind()
	query := bson.M{"topicId": req.TopicId}
	cur, err := col.Find(
		findCtx, query, &options.FindOptions{Skip: &start},
	)
	if err != nil {
		return
	}

	chstream, err := col.Watch(
		stream.Context(),
		bson.M{"$match": query},
	)
	defer chstream.Close(stream.Context())
	return
}
