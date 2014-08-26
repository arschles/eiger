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
	"os"
)

func serve(wsConn *websocket.Conn, dclient *docker.Client, diedCh chan<- error) {
	handlers := NewHandlers(dclient)
	server := rpc.NewServer()
	server.Register(handlers)
	serverCodec := jsonrpc.NewServerCodec(wsConn)
	server.ServeCodec(serverCodec)
	diedCh <- fmt.Errorf("server stopped serving")
}

func heartbeat(wsConn *websocket.Conn, interval time.Duration, diedCh chan<- error) {
	hostname, err := os.Hostname()
	if err != nil {
		diedCh <- err
		return
	}
	clientCodec := jsonrpc.NewClientCodec(wsConn)
	client := rpc.NewClientWithCodec(clientCodec)
	for {
		rep := 1
		err := client.Call("Handlers.Heartbeat", hostname, &rep)
		if err != nil {
			util.LogWarnf("(error heartbeating) %s", err)
		} else if rep != 0 {
			util.LogWarnf("(error heartbeating) expected return code was %d, not 0", rep)
		}
		time.Sleep(interval)
	}
	diedCh <- fmt.Errorf("heartbeat loop stopped")
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

	serveDied := make(chan error)
	go serve(wsConn, dclient, serveDied)
	heartbeatDied := make(chan error)
	go func() {
		heartbeat(wsConn, hbInterval, heartbeatDied)
	}()

	for {
		select {
		case err := <-serveDied:
			log.Fatalf("(rpc server) %s", err)
		case err := <-heartbeatDied:
			util.LogWarnf("(heartbeat loop) %s", err)
			go heartbeat(wsConn, hbInterval, heartbeatDied)
		}
	}
}
