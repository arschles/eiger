package pubsub

import (
  "log"
)

type LoggingPublisher struct {
}

func (n LoggingPublisher) Publish(p *Payload) {
  log.Printf("(%s) %s", p.Topic, p.Data)
}
