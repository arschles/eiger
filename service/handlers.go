package main

import (
    "log"
    "time"
)


type Handlers struct {
    set *AgentSet
    heartbeat time.Duration
}

func NewHandlers(set *AgentSet, hb time.Duration) *Handlers {
    return &Handlers{set, hb}
}

func (h *Handlers) watchAgent(agent Agent) {
    for {
        //sleep for 2x the heartbeat duration, make
        //sure we got at least 1 heartbeat afterward
    }
}

func (h *Handlers) Heartbeat(host string, rep *int) error {
    *rep = 0
    agent := NewAgent(host)
    if !h.set.Add(agent) {
        log.Printf("added agent %s to alive set", agent)
        go h.watchAgent(agent)
    }
    return nil
}
