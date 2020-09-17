package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	pb "github.com/kubnix/skipper/protos/event"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:42113", grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		panic(err)
	}

	client := pb.NewEventSvcClient(conn)
	stream, err := client.Watch(context.Background(), &pb.WatchRequest{
		EventType:    pb.EventType_CREATED,
		ResourceType: pb.ResourceType_WORKFLOW,
	})
	for event, err := stream.Recv(); err == nil; event, err = stream.Recv() {
		// log.Info(event.Data)
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, event.Data, "", "\t")
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Println(string(prettyJSON.Bytes()))
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
