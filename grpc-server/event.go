package grpcserver

import (
	"github.com/kubnix/skipper/listener"
	pb "github.com/kubnix/skipper/protos/event"
)

func (s *server) Watch(req *pb.WatchRequest, stream pb.EventSvc_WatchServer) error {
	err := listener.StartListener(func(ev *pb.Event) error {
		return stream.Send(ev)
	})
	return err
}
