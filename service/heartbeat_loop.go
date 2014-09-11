package main

import (
	"github.com/arschles/eiger/lib/util"
	"time"
)

//agentWatcher runs a watch loop on agent, expecting a heartbeat on ticker every
//interval. if it doesn't get a heartbeat in interval (plus a built in grace
//period), the loop stops and sends on removed
func agentWatcher(agent Agent,
	interval time.Duration,
	ticker <-chan interface{},
	removed chan<- interface{}) {

	for {
		start := time.Now()
		select {
		case <-ticker:
			if time.Since(start) > interval*4 {
				util.LogWarnf("(late heartbeat) removing agent %s from alive set", agent)
				break
			}
		case <-time.After(interval * 2):
			util.LogWarnf("(heartbeat timeout) removing agent %s from alive set", agent)
			break
		}
	}
	removed <- struct{}{}
}
