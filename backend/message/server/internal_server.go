package server

import (
	prvRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v4"
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
	var numDoc int64
	var topics []*prvRPC.StatsResponse

	for {
		req, err = srv.Recv()
		if err != nil {
			break
		}
		var topicID pr.ObjectID
		topicID, err = pr.ObjectIDFromHex(req.GetTopicId())
		if err != nil {
			break
		}
		col := me.collection()
		numDoc, err = col.CountDocuments(srv.Context(), bson.M{"topicid": topicID})
		if err != nil {
			break
		}
		resp := &prvRPC.StatsResponse{
			TopicId: topicID.Hex(),
			NumMsgs: numDoc,
		}
		err = srv.Send(resp)
		if err != nil {
			break
		}
		topics <- &statsSubscriber{
			TopicID: topicID,
			Resp:    resp,
		}
		sub, err := me.subscribe(topicID.Hex(), topics[topicID].MsgCh)
		if err != nil {
			break
		}
		defer sub.Unsubscribe()

		go func() {
			msgCh := make(chan *nats.Msg)
			defer close(msgCh)
			sub, err := me.subscribe(topicID.Hex(), msgCh)
			if err != nil {
				return
			}
			defer sub.Unsubscribe()
			for {
				select {
				case rec := <-msgCh:
					var msg rpc.Message
					err = msgpack.Unmarshal(rec.Data, &msg)
					if err != nil {
						return
					}
					resp.NumMsgs++
					if msg.GetBump() {
						resp.LastBump = msg.GetPostTime()
					}
					srv.Send(resp)
				case <-srv.Context().Done():
					return
				}
			}
		}()
	}
	return
}
