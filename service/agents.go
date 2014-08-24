package main

import (
	"io"
	"net/url"
	"sync"
	"time"
)

//Agent is the full representation of an agent, including the io.Writer
//that can be used to communicate with the agent
type Agent struct {
	Origin     *url.URL
	ReadWriter io.ReadWriter
}

func NewAgent(o *url.URL, rw io.ReadWriter) *Agent {
	return &Agent{o, rw}
}

func (a *Agent) String() string {
	return a.Origin.String()
}

//Agents represents a set of agents, each of which must have a heartbeat on its
//Writer. When the heartbeat fails, the agent is removed from the set
type Agents struct {
	agents map[Agent]bool
	mutex  sync.RWMutex
	hb     time.Duration
}

func NewAgents(agents *[]Agent, hb time.Duration) *Agents {
	m := map[Agent]bool{}
	for _, a := range *agents {
		m[a] = true
	}
	return &Agents{
		agents: m,
		hb:     hb,
	}
}

//Add adds an Agent to this set
func (a *Agents) Add(agnt Agent) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	_, ok := a.agents[agnt]
	if !ok {
		return false
	}
	a.agents[agnt] = true
	return true
}

//Remove removes the given agent from the internal set
func (a *Agents) Remove(agnt Agent) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	_, ok := a.agents[agnt]
	if !ok {
		return false
	}
	delete(a.agents, agnt)
	return true
}
