package server

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/hiroaki-yamamoto/real/backend/config"
	"github.com/hiroaki-yamamoto/real/backend/rpc"
	"github.com/hiroaki-yamamoto/real/backend/validation"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/peer"
)

// Server implements MessageServiceServer interface.
type Server struct {
	Setting  *config.Config
	Database *mongo.Database
	Broker   *nats.Conn
}

func (me *Server) getCollection() *mongo.Collection {
	return me.Database.Collection(srvName)
}

func (me *Server) getBrokerSubject(topicID string) string {
	return srvName + "/" + topicID
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

	col := me.getCollection()
	query := bson.M{"topicid": topicID}
	findCur, err := col.Find(
		stream.Context(), query,
		&options.FindOptions{
			Skip: &start,
			Sort: bson.M{"posttime": 1},
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
		err = stream.Send(model.ToRPCMsg(true))
		if err != nil {
			return
		}
	}

	msgCh := make(chan *nats.Msg)
	defer close(msgCh)
	chSub, err := me.Broker.ChanSubscribe(
		me.getBrokerSubject(req.TopicId), msgCh,
	)
	if err != nil {
		return
	}
	defer chSub.Unsubscribe()
	go me.Broker.Publish("status/"+srvName+"/subscribe", []byte("ready"))
	for {
		select {
		case msg := <-msgCh:
			var model rpc.Message
			if err = msgpack.Unmarshal(msg.Data, &model); err != nil {
				return
			}
			// model.SenderName = html.EscapeString(model.SenderName)
			// model.Message = html.EscapeString(model.Message)
			stream.Send(&model)
			break
		case <-stream.Context().Done():
			return
		}
	}
}

// Post recirds the message, broadcast it, and returns Status structure
func (me *Server) Post(
	ctx context.Context,
	req *rpc.PostRequest,
) (status *rpc.Status, err error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		err = errors.New("Your IP Address is not contained in the current context")
		return
	}
	remoteIP, _, err := net.SplitHostPort(peer.Addr.String())
	if err != nil {
		return
	}
	vld, err := validation.New(ctx, me.Setting.Recaptcha)
	if err != nil {
		return
	}
	topicID, err := primitive.ObjectIDFromHex(req.GetTopicId())
	if err != nil {
		return
	}
	model := &Model{
		ID:         primitive.NewObjectID(),
		TopicID:    topicID,
		SenderName: req.GetName(),
		PostTime:   time.Now().UTC(),
		Message:    req.GetMessage(),
		Host:       remoteIP,
		Bump:       req.GetBump(),
	}
	err = vld.Struct(model)
	if err != nil {
		return
	}
	err = model.Store(ctx, me.getCollection())
	if err != nil {
		return
	}
	msg, err := msgpack.Marshal(model.ToRPCMsg(false))
	if err != nil {
		return
	}
	err = me.Broker.Publish(me.getBrokerSubject(req.GetTopicId()), msg)
	if err != nil {
		return
	}
	status = &rpc.Status{
		Id: model.ID.Hex(),
	}
	return
}
