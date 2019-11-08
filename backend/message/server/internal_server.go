package server

import (
	"errors"

	prvRPC "github.com/hiroaki-yamamoto/real/backend/message/rpc"
)

// InternalServer is a server to provide internal information like stats info.
type InternalServer struct {
}

// Stats generates statistics report of the specified message
func (me *InternalServer) Stats(
	srv prvRPC.MessageStats_StatsServer,
) error {
	return errors.New("Not Implemented Yet")
}
