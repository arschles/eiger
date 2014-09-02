package main

import (
    "log"
    "time"
    "github.com/arschles/eiger/lib/util"
)


type AgentHandlers struct {
    set *AgentSet
    heartbeat time.Duration
    //channel to signal to start watching an agent
    watchCh chan Agent
    tickCh chan Agent
}

func NewAgentHandlers(set *AgentSet, hb time.Duration) *AgentHandlers {
    h := Handlers {
        set: set,
        heartbeat: hb,
        watchCh: make(chan Agent),
        tickCh: make(chan Agent),
    }
    go h.tickLoop()
    return &h
}

func (h *AgentHandlers) tickLoop() {

    watchAgent := func(agent Agent, ticker <-chan bool) {
        for {
            start := time.Now()
            select {
            case <-ticker:
                if time.Since(start) > h.heartbeat * 4 {
                    util.LogWarnf("(late heartbeat) removing agent %s from alive set", agent)
                    h.set.Remove(agent)
                    return
                }
            case <-time.After(h.heartbeat * 2):
                util.LogWarnf("(heartbeat timeout) removing agent %s from alive set", agent)
                return
            }
        }
    }

    tickers := map[Agent]chan bool{}
    for {
        select {
        case agent := <-h.watchCh:
            tickerCh := make(chan bool)
            tickers[agent] = tickerCh
            go watchAgent(agent, tickerCh)
        case agent := <-h.tickCh:
            tickerCh, ok := tickers[agent]
            if !ok {
                util.LogWarnf("could not find ticker channel for agent %s", agent)
                continue
            }
            go func() {
                tickerCh <- true
            }()
        }
    }
}

func (h *AgentHandlers) Heartbeat(host string, rep *int) error {
    *rep = 0
    agent := NewAgent(host)
    added := h.set.Add(agent)
    if added {
        log.Printf("added agent %s to alive set", agent)
        h.watchCh <- agent
    } else {
        h.tickCh <- agent
    }
    return nil
}
