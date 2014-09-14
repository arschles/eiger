package pubsub

import (
	"log"
)

const (
	HeartbeatTopic         = "heartbeat"
	HeartbeatErrorTopic    = "heartbeat_error"
	LateHeartbeatTopic     = "heartbeat_late"
	DockerEventsTopic      = "dockerevt"
	DockerEventsErrorTopic = "dockerevt_error"
)

type Payload struct {
	Topic string
	Data  interface{}
}

//Publisher publishes payloads to somewhere
type Publisher interface {
	Publish(p *Payload)
}

//Subscriber gets payloads that have a specific topic
type Subscriber interface {
	Subscribe(topic string) <-chan *Payload
}

type LoggingPublisher struct {
}

func (n LoggingPublisher) Publish(p *Payload) {
	log.Printf("(%s) %s", p.Topic, p.Data)
}
