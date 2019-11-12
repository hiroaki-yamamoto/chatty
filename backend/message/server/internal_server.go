package server

import (
	prvRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	pr "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InternalServer is a server to provide internal information like stats info.
type InternalServer struct {
	DB     *mongo.Database
	Broker *nats.Conn
}

func (me *InternalServer) collection() *mongo.Collection {
	return me.DB.Collection(srvName)
}

func (me *InternalServer) subscribe(
	topicID string,
	ch chan *nats.Msg,
) (*nats.Subscription, error) {
	return me.Broker.ChanSubscribe(srvName+"/"+topicID, ch)
}

// Stats generates statistics report of the specified message
func (me *InternalServer) Stats(
	srv prvRPC.MessageStats_StatsServer,
) (err error) {
	var req *prvRPC.StatsRequest
	var topicID pr.ObjectID
	var numDoc int64
	req, err = srv.Recv()
	if err != nil {
		return nil
	}
	topicID, err = pr.ObjectIDFromHex(req.GetTopicId())
	if err != nil {
		return nil
	}
	col := me.collection()
	numDoc, err = col.CountDocuments(srv.Context(), bson.M{"topicid": topicID})
	// col.Find(srv.Context(), bson.M{"topicid": topicID, "dump": true})
	if err != nil {
		return
	}
	resp := &prvRPC.StatsResponse{
		TopicId: topicID.Hex(),
		NumMsgs: numDoc,
	}
	err = srv.Send(resp)
	if err != nil {
		return
	}
	msgCh := make(chan *nats.Msg)
	defer close(msgCh)
	sub, err := me.subscribe(topicID.Hex(), msgCh)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()
	for {
		select {
		case <-msgCh:
			// Decode the data
		case <-srv.Context().Done():
			return
		}
	}
}
