package main

import (
	"sync"
)

//Agent is the full representation of an agent, including the io.Writer
//that can be used to communicate with the agent
type Agent string

func NewAgent(name string) Agent {
	return Agent(name)
}

//Agents represents a set of agents, each of which must have a heartbeat on its
//Writer. When the heartbeat fails, the agent is removed from the set
type AgentSet struct {
	set map[Agent]bool
	mutex  sync.RWMutex //protects set
}

func NewAgentSet(agents *[]Agent) *AgentSet {
	m := map[Agent]bool{}
	for _, a := range *agents {
		m[a] = true
	}
	return &AgentSet{
		set: m,
	}
}

//Add adds an Agent to this set
func (a *AgentSet) Add(agnt Agent) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	_, ok := a.set[agnt]
	if !ok {
		return false
	}
	a.set[agnt] = true
	return true
}

//Remove removes the given agent from the internal set
func (a *AgentSet) Remove(agnt Agent) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	_, ok := a.set[agnt]
	if !ok {
		return false
	}
	delete(a.set, agnt)
	return true
}
