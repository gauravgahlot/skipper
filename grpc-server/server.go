package grpcserver

import (
	"context"
	"net"

	"github.com/kubnix/skipper/db"
	"github.com/kubnix/skipper/protos/event"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type server struct {
	db db.Database
}

// SetupGRPC setup and return a gRPC server
func SetupGRPC(ctx context.Context, errCh chan<- error) {
	server := &server{
		db: db.Connect(),
	}

	s := grpc.NewServer(grpc.EmptyServerOption{})
	event.RegisterEventSvcServer(s, server)

	go func() {
		lis, err := net.Listen("tcp", ":42113")
		if err != nil {
			err = errors.Wrap(err, "failed to listen")
			log.Error(err)
			panic(err)
		}
		errCh <- s.Serve(lis)
	}()

	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()
}
