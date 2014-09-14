package pubsub

import (
  "sync"
)

type InMemPublisherSubscriber struct {
  mut sync.RWMutex
  subscribers map[string] []chan<-*Payload
}

func (i InMemPublisherSubscriber) Publish(p *Payload) {
  i.mut.RLock()
  defer i.mut.RUnlock()
  chans, ok := i.subscribers[p.Topic]
  if !ok {
    return
  }
  for _, ch := range chans {
    go func () {
      ch <- p
    }()
  }
}

func (i *InMemPublisherSubscriber) Subscribe(topic string) <-chan *Payload{
  i.mut.Lock()
  defer i.mut.Unlock()
  newCh := make(chan *Payload)
  lst, ok := i.subscribers[topic]
  if !ok {
    lst = []chan<- *Payload {}
  }
  lst = append(lst, newCh)
  i.subscribers[topic] = lst
  return newCh
}
