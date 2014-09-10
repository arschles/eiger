package main

import (
	"fmt"
	"github.com/arschles/eiger/lib/util"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"log"
	"time"
	"code.google.com/p/go.net/websocket"
)

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
	go rpcLoop(wsConn, dclient, serveDied)
	log.Printf("started RPC server")

	heartbeatDied := make(chan error)
	go heartbeatLoop(wsConn, hbInterval, heartbeatDied)
	log.Printf("started heartbeat loop")

	logsDied := make(chan error)
	go logsLoop(wsConn, logsDied)
	log.Printf("started logs loop")

	for {
		select {
		case err := <-serveDied:
			log.Fatalf("(rpc server) %s", err)
		case err := <-heartbeatDied:
			util.LogWarnf("(heartbeat loop) %s", err)
			go heartbeatLoop(wsConn, hbInterval, heartbeatDied)
		case err := <- logsDied:
			util.LogWarnf("(logs loop) %s", err)
			go logsLoop(wsConn, logsDied)
		}
	}
}
