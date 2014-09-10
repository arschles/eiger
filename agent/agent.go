package main

import (
	"fmt"
	"github.com/arschles/eiger/lib/util"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"log"
	"time"
	"net/rpc"
	"net/rpc/jsonrpc"
	"code.google.com/p/go.net/websocket"
)

func serve(wsConn *websocket.Conn, dclient *docker.Client, diedCh chan<- error) {
	handlers := NewHandlers(dclient)
	server := rpc.NewServer()
	server.Register(handlers)
	serverCodec := jsonrpc.NewServerCodec(wsConn)
	for {
		server.ServeCodec(serverCodec)
	}
	diedCh <- fmt.Errorf("server stopped serving")
}

func agent(c *cli.Context) {
	dclient, err := docker.NewClient(c.String("dockerhost"))
	if err != nil {
		log.Fatalf("(docker connection) %s", err)
	}

	host := c.String("host")
	port := c.Int("port")
	hbInterval := time.Duration(c.Int("heartbeat")) * time.Millisecond
	heartbeatUrl := fmt.Sprintf("ws://%s:%d/heartbeat", host, port)
	origin := fmt.Sprintf("http://%s/", host)

	log.Printf("dialing %s", heartbeatUrl)
	wsConn, err := websocket.Dial(heartbeatUrl, "", origin)
	if err != nil {
		log.Fatalf("(websocket connection) %s", err)
	}

	serveDied := make(chan error)
	go serve(wsConn, dclient, serveDied)
	log.Printf("started RPC server")

	heartbeatDied := make(chan error)
	go heartbeatLoop(wsConn, hbInterval, heartbeatDied)
	log.Printf("started heartbeat loop")


	for {
		select {
		case err := <-serveDied:
			log.Fatalf("(rpc server) %s", err)
		case err := <-heartbeatDied:
			util.LogWarnf("(heartbeat loop) %s", err)
			go heartbeatLoop(wsConn, hbInterval, heartbeatDied)
		}
	}
}
