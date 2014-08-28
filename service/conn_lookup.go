package main

import (
    "code.google.com/p/go.net/websocket"
    "sync"
)

type ConnLookup struct {
    unmatched map[Host]*websocket.Conn
    unmatchedLock sync.RWMutex
    matched map[Agent]*websocket.Conn
}

func NewConnLookup() *ConnLookup {
    return &ConnLookup{}
}

func (c *ConnLookup) AddUnmatched(conn *websocket.Conn) {
    c.unmatchedLock.Lock()
    defer c.unmatchedLock.Unlock()
    host := conn.Config().Origin.Host
    c.unmatched[host] = conn
}

func (c *ConnLookup) Match(host Host) *websocket.Conn {
    c.unmatchedLock.RLock()
    
}
