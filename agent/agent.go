package main

import (
	"fmt"
	"github.com/arschles/eiger/lib/util"
	"github.com/arschles/eiger/lib/ws"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"log"
	"time"
	"net/rpc"
	"net/rpc/jsonrpc"
	"code.google.com/p/go.net/websocket"
)

func serve(wsConn *websocket.Conn, dclient *docker.Client, diedCh chan<- bool) {
	handlers := NewHandlers(dclient)
	server := rpc.NewServer()
	server.Register(handlers)
	serverCodec := jsonrpc.NewServerCodec(wsConn)
	server.ServeCodec(serverCodec)
	diedCh <- true
}

func heartbeat(wsConn *websocket.Conn, interval time.Duration, diedCh chan<- bool) {
	clientCodec := jsonrpc.NewClientCodec(wsConn)
	client := rpc.NewClientWithCodec(clientCodec)
	for {
		res := struct{}{}
		rep := struct{}{}
		err := client.Call("Handlers.Heartbeat", res, &rep)
		if err != nil {
			util.LogWarnf("error heartbeating: %s", err)
		}
	}
	diedCh <- true
}

func agent(c *cli.Context) {
	dclient, err := docker.NewClient(c.String("dockerhost"))
	if err != nil {
		log.Fatal(err)
	}

	host := c.String("host")
	port := c.Int("port")
	hbInterval := time.Duration(c.Int("heartbeat")) * time.Millisecond
	socketUrl := fmt.Sprintf("ws://%s:%d/socket", host, port)
	origin := fmt.Sprintf("http://%s/", host)

	log.Printf("dialing %s", socketUrl)
	wsConn := ws.MustDial(socketUrl, origin)

	serveDied := make(chan bool)
	go serve(wsConn, dclient, serveDied)
	heartbeatDied := make(chan bool)
	go heartbeat(wsConn, hbInterval, heartbeatDied)

	for {
		select {
		case <-serveDied:
			log.Fatal("agent server died")
		case <-heartbeatDied:
			util.LogWarnf("heartbeat producer died, restarting")
			go heartbeat(wsConn, hbInterval, heartbeatDied)
		}
	}
}
