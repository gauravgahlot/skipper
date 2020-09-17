package listener

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/kubnix/skipper/protos/event"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	eventsChannel = "events_channel"
	connInfo      = "dbname=tinkerbell user=tinkerbell password=tinkerbell sslmode=disable"
)

// StartListener creates a new dedicated connection for LISTEN/NOTIFY
// and starts listening for events.
func StartListener(fn func(e *pb.Event) error) error {
	_, err := sql.Open("postgres", "")
	if err != nil {
		log.Error(err)
	}

	listener := pq.NewListener(connInfo, 5*time.Second, 15*time.Second, errorHandler)
	err = listener.Listen(eventsChannel)
	if err != nil {
		log.Error(err)
	}

	log.Info("starting listener")
	for {
		err := waitForNotification(listener, fn)
		if err != nil {
			return err
		}
	}
}

func waitForNotification(l *pq.Listener, fn func(e *pb.Event) error) error {
	for {
		n := <-l.Notify
		// log.Info("Received data from channel [", n.Channel, "] :")

		var ev event
		err := json.Unmarshal([]byte(n.Extra), &ev)
		if err != nil {
			log.Error(err)
			return err
		}

		d, _ := json.Marshal(ev.Data)
		cAt, _ := ptypes.TimestampProto(*ev.CreatedAt)
		newEvent := pb.Event{
			Id:           ev.ID,
			ResourceId:   ev.ResourceID,
			ResourceType: pb.ResourceType_WORKFLOW, // get resource type based on ev.ResourceType int value
			EventType:    pb.EventType_CREATED,     // get resource type based on ev.EventType int value
			Data:         d,
			CreatedAt:    cAt,
		}
		err = fn(&newEvent)
		if err != nil {
			log.Error(err)
			return err
		}
	}
}

func errorHandler(ev pq.ListenerEventType, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

type event struct {
	ID           string      `json:"id,omitempty"`
	ResourceID   string      `json:"resource_id,omitempty"`
	ResourceType int32       `json:"resource_type,omitempty"`
	EventType    int32       `json:"event_type,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	CreatedAt    *time.Time  `json:"created_at,omitempty"`
}
