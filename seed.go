package main

import (
	"encoding/json"
	"time"

	"github.com/kubnix/skipper/db"
	pb "github.com/kubnix/skipper/protos/event"
	"github.com/prometheus/common/log"
)

var workflowID = []string{
	"c04a8229-0856-415f-b750-a7401ad447b0",
	"1f2b8e74-508a-4b03-8a3f-48d84c58d916",
	"60f1303a-dc1c-4353-a2a0-cdeb52cf25ef",
}

func sleepAndRun() {
	d := db.Connect()
	// time.Sleep(10 * time.Second)
	log.Info("generating dummy events...")

	for _, wfID := range workflowID {
		time.Sleep(10 * time.Second)
		d.AddEvent(
			wfID,
			pb.ResourceType_value[pb.ResourceType_WORKFLOW.String()],
			pb.EventType_value[pb.EventType_CREATED.String()],
			workflow(wfID),
		)
	}
}

// Workflow represents a workflow to be executed
type Workflow struct {
	Version       string `json:"version"`
	Name          string `json:"name"`
	ID            string `json:"id"`
	GlobalTimeout int    `json:"global_timeout"`
	Tasks         []Task `json:"tasks"`
}

// Task represents a task to be executed as part of a workflow
type Task struct {
	Name        string            `json:"name"`
	WorkerAddr  string            `json:"worker"`
	Actions     []Action          `json:"actions"`
	Volumes     []string          `json:"volumes"`
	Environment map[string]string `json:"environment"`
}

// Action is the basic executional unit for a workflow
type Action struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Timeout     int64             `json:"timeout"`
	Command     []string          `json:"command"`
	OnTimeout   []string          `json:"on-timeout"`
	OnFailure   []string          `json:"on-failure"`
	Volumes     []string          `json:"volumes,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

func workflow(wfID string) []byte {
	wf := &Workflow{
		ID:            wfID,
		GlobalTimeout: 900,
		Name:          "ubuntu-provisioning",
		Tasks: []Task{
			{
				Name:       "pre-installation",
				WorkerAddr: "08:00:27:00:00:01",
				Environment: map[string]string{
					"MIRROR_HOST": "192.168.1.2",
				},
				Volumes: []string{
					"/dev:/dev",
					"/dev/console:/dev/console",
					"/lib/firmware:/lib/firmware:ro",
				},
				Actions: []Action{
					{
						Name:    "disk-partition",
						Image:   "disk-partition",
						Timeout: 300,
						Volumes: []string{
							"/statedir:/statedir",
						},
					},
					{
						Name:    "install-root-fs",
						Image:   "install-root-fs",
						Timeout: 600,
					},
				},
			},
		},
	}
	data, err := json.Marshal(wf)
	if err != nil {
		log.Error(err)
		return nil
	}
	return data
}
