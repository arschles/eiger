package main

import (
    "time"
    "github.com/arschles/eiger/lib/util"
)

type HeartbeatLoop struct {
    lookup *AgentLookup
    hbDur time.Duration
    notifyCh chan Agent
}

func NewHeartbeatLoop(l *AgentLookup, hbDur time.Duration) *HeartbeatLoop {
    notifyCh := make(chan Agent)
    loop := HeartbeatLoop{l, hbDur, notifyCh}

    watchCh := make(chan Agent)
    tickCh := make(chan Agent)
    go loop.run(watchCh, tickCh)
    
    return &loop
}

//Notify tells the heartbeat loop that an agent has either heartbeated
//or has been added
func (h *HeartbeatLoop) Notify(a Agent) {
    h.notifyCh <- a
}

func (h *HeartbeatLoop) agentWatcher(agent Agent, ticker <-chan bool) {
  for {
      start := time.Now()
      select {
      case <-ticker:
          if time.Since(start) > h.hbDur * 4 {
              util.LogWarnf("(late heartbeat) removing agent %s from alive set", agent)
              h.lookup.Remove(agent)
              return
          }
      case <-time.After(h.hbDur * 2):
          util.LogWarnf("(heartbeat timeout) removing agent %s from alive set", agent)
          return
      }
  }
}

func (h *HeartbeatLoop) run(watchCh chan Agent, tickCh chan Agent) {
    tickers := map[Agent]chan bool{}
    for {
        select {
        case agent := <-watchCh:
            tickerCh := make(chan bool)
            tickers[agent] = tickerCh
            go h.agentWatcher(agent, tickerCh)
        case agent := <-tickCh:
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
