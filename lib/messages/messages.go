package messages

import (
	"time"
)

//HeartbeatMessage represents the JSON structure that's sent over the wire
type Heartbeat struct {
	Hostname string    `json:"hostname"`
	SendTime time.Time `json:"time"`
}

type DockerEvent struct {
}

type DockerLog struct {
	Container string `json:"container"`
	Out       string `json:"output"`
	Err       string `json:"error"`
}

type DockerStatus struct {
}
