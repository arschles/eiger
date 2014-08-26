package main

type Handlers struct {
    set *AgentSet

}

func NewHandlers(set *AgentSet) *Handlers {
    return &Handlers{set}
}

func (h *Handlers) Heartbeat(_ struct{}, rep *int) error {
    *rep = 0
    return nil
}
