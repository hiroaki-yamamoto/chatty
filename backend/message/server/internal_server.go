package server

import (
	prvRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
	"go.mongodb.org/mongo-driver/bson"
	pr "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InternalServer is a server to provide internal information like stats info.
type InternalServer struct {
	DB *mongo.Database
}

func (me *InternalServer) collection() *mongo.Collection {
	return me.DB.Collection(srvName)
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
	for {
		select {
		case <-srv.Context().Done():
			return
		}
	}
}
