package server

import (
	"sync"

	"github.com/golang/protobuf/ptypes/timestamp"
	prvRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v4"
	"go.mongodb.org/mongo-driver/bson"
	pr "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	statsStore := make(map[pr.ObjectID]*prvRPC.StatsResponse)

	var subscriptions []*nats.Subscription
	defer func() {
		for _, sub := range subscriptions {
			sub.Unsubscribe()
		}
	}()

	msgCh := make(chan *nats.Msg, 1024)
	defer close(msgCh)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
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
			var lastBumpDoc Model
			lastBumpDocCur, err := col.Find(
				srv.Context(),
				bson.M{"bump": true, "topicid": topicID},
				options.Find().SetSort(
					bson.M{"posttime": -1},
				).SetLimit(1),
			)
			if err != nil {
				break
			}
			for lastBumpDocCur.Next(srv.Context()) {
				lastBumpDocCur.Decode(&lastBumpDoc)
			}
			resp := &prvRPC.StatsResponse{
				TopicId: topicID.Hex(),
				NumMsgs: numDoc,
				LastBump: &timestamp.Timestamp{
					Seconds: lastBumpDoc.PostTime.Unix(),
					Nanos:   int32(lastBumpDoc.PostTime.Nanosecond()),
				},
			}
			err = srv.Send(resp)
			if err != nil {
				break
			}
			sub, err := me.subscribe(topicID.Hex(), msgCh)
			if err != nil {
				break
			}
			subscriptions = append(subscriptions, sub)
			statsStore[topicID] = resp
			select {
			case <-srv.Context().Done():
				return
			default:
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case rec := <-msgCh:
				var msg rpc.Message
				err = msgpack.Unmarshal(rec.Data, &msg)
				if err != nil {
					return
				}
				topicID, err := pr.ObjectIDFromHex(msg.TopicId)
				if err != nil {
					return
				}
				resp := statsStore[topicID]
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

	wg.Wait()
	return
}
