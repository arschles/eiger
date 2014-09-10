package main

import (
	"io"
	"sync"
)

//Agent is the full representation of an agent, including the io.Writer
//that can be used to communicate with the agent
type Agent struct {
	Hostname string
	Conn     io.Writer
}

func NewAgent(hostname string, conn io.Writer) *Agent {
	return &Agent{hostname, conn}
}

func (a Agent) String() string {
	return a.Hostname
}

//AgentLookup represents a set of agents, each of which must have a heartbeat on its
//Writer. When the heartbeat fails, the agent is removed from the set
type AgentLookup struct {
	m     map[string]Agent
	mutex sync.RWMutex //protects set
}

func NewAgentLookup(agents *[]Agent) *AgentLookup {
	m := map[string]Agent{}
	for _, agent := range *agents {
		m[agent.Hostname] = agent
	}
	return &AgentLookup{
		m: m,
	}
}

//Add adds an Agent to this set
func (a *AgentLookup) GetOrAdd(agent Agent) *Agent {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	existing, ok := a.m[agent.Hostname]
	if !ok {
		a.m[agent.Hostname] = agent
		return &agent
	}
	return &existing
}

//Remove removes the given agent from the internal set
func (a *AgentLookup) Remove(agnt Agent) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	_, ok := a.m[agnt.Hostname]
	if !ok {
		return false
	}
	delete(a.m, agnt.Hostname)
	return true
}
