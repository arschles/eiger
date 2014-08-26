package main

import (
    "code.google.com/p/go.net/websocket"
)

type ConnLookup struct {

}

func NewConnLookup() *ConnLookup {
    return &ConnLookup{}
}

func (c *ConnLookup) AddUnmatched(conn *websocket.Conn) {

}

func (c *ConnLookup) Match(hostname string) *websocket.Conn {
    
}
