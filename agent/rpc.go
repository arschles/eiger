package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func rpcLoop(wsConn *websocket.Conn, dclient *docker.Client, diedCh chan<- error) {
	handlers := NewHandlers(dclient)
	server := rpc.NewServer()
	server.Register(handlers)
	serverCodec := jsonrpc.NewServerCodec(wsConn)
	for {
		server.ServeCodec(serverCodec)
	}
}
