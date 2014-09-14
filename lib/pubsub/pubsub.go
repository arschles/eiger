package pubsub

import (
  "log"
)

type Payload interface{
  Topic() string
  Data() interface{}
}

//Publisher publishes payloads to somewhere
type Publisher interface {
  Publish(p Payload)
}

//Subscriber gets payloads that have a specific topic
type Subscriber interface {
  Subscribe(topic string) <-chan Payload
}

type LoggingPublisher struct {

}
func (n *NoopPublisher) Publish(p Payload) {
  log.Printf("publishing (%s) %s" p.Topic(), p.Data())
}
