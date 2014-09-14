package main

import (
	"github.com/arschles/eiger/lib/util"
	"time"
)

//the multiplier for heartbeat duration. used to determine when a heartbeat
//has timed out
const HB_TIMEOUT_MULTIPLIER = 10

type HeartbeatLoop struct {
	lookup   *AgentLookup
	hbDur    time.Duration
	notifyCh chan Agent
}

func NewHeartbeatLoop(l *AgentLookup, hbDur time.Duration) *HeartbeatLoop {
	notifyCh := make(chan Agent)
	loop := HeartbeatLoop{l, hbDur, notifyCh}

	go loop.run()

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
			if time.Since(start) > h.hbDur*4 {
				util.LogWarnf("(late heartbeat) removing agent %s from alive set", agent)
				h.lookup.Remove(agent)
				return
			}
		case <-time.After(h.hbDur * HB_TIMEOUT_MULTIPLIER):
			util.LogWarnf("(heartbeat timeout) removing agent %s from alive set", agent)
			h.lookup.Remove(agent)
			return
		}
	}
}

func (h *HeartbeatLoop) run() {
	tickers := map[Agent]chan bool{}
	for {
		select {
		case agent := <-h.notifyCh:
			tickerCh, ok := tickers[agent]
			if !ok {
				t := make(chan bool)
				tickers[agent] = t
				tickerCh = t
			}
			go h.agentWatcher(agent, tickerCh)
			go func() {
				tickerCh <- true
			}()
		}
	}
}
