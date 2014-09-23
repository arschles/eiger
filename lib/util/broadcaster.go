package util

import (
  "sync"
)

type Broadcaster struct {
  source <-chan interface{}
  targets []chan<- interface{}
  mut sync.RWMutex
}

func NewBroadcaster(source <-chan interface{}) *Broadcaster {
  targets := []chan<- interface{}{}
  var mut sync.RWMutex
  b := Broadcaster {
    source: source,
    targets: targets,
    mut: mut,
  }

  go func() {
    for {
      value := <-b.source
      b.mut.RLock()
      for _, target := range b.targets {
        go func(target chan<-interface{}) {
          target <- value
        }(target)
      }
      b.mut.RUnlock()
    }
  }()

  return &b
}

func (b *Broadcaster) NewChan() <-chan interface{} {
  b.mut.Lock()
  defer b.mut.Unlock()
  newCh := make(chan interface{})
  b.targets = append(b.targets, newCh)
  return newCh
}
