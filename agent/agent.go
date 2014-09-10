package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/arschles/eiger/lib/util"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"log"
	"time"
)

func dialOrDie(url string, origin string) *websocket.Conn {
	log.Printf("dialing %s", url)
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatalf("(websocket connection) %s", err)
	}
	return conn
}

func startRpc(origin string, host string, port int, dclient *docker.Client) <-chan error {
	rpcConn := dialOrDie(fmt.Sprintf("ws://%s:%d/rpc", host, port), origin)
	serveDied := make(chan error)
	go rpcLoop(rpcConn, dclient, serveDied)
	log.Printf("started RPC server")
	return serveDied
}

func startHb(origin string, host string, port int, interval time.Duration) <-chan error {
	hbConn := dialOrDie(fmt.Sprintf("ws://%s:%d/heartbeat", host, port), origin)
	heartbeatDied := make(chan error)
	go heartbeatLoop(hbConn, interval, heartbeatDied)
	log.Printf("started heartbeat loop")
	return heartbeatDied
}

func startDocker(origin string, host string, port int, client *docker.Client) <-chan error {
	dockerConn := dialOrDie(fmt.Sprintf("ws://%s:%d/docker", host, port), origin)
	dockerCh := make(chan error)
	go dockerLoop(dockerConn, client, dockerCh)
	log.Printf("started docker loop")
	return dockerCh
}

func agent(c *cli.Context) {
	dclient, err := docker.NewClient(c.String("dockerhost"))
	if err != nil {
		log.Fatalf("(docker connection) %s", err)
	}

	host := c.String("host")
	port := c.Int("port")
	hbInterval := time.Duration(c.Int("heartbeat")) * time.Millisecond
	origin := fmt.Sprintf("http://%s/", host)

	rpcDied := startRpc(origin, host, port, dclient)
	hbDied := startHb(origin, host, port, hbInterval)
	dockerDied := startDocker(origin, host, port, dclient)

	for {
		select {
		case err := <-rpcDied:
			log.Fatalf("(rpc server) %s", err)
		case err := <-hbDied:
			util.LogWarnf("(heartbeat loop) %s", err)
			hbDied = startHb(origin, host, port, hbInterval)
		case err := <-dockerDied:
			util.LogWarnf("(docker loop) %s", err)
			dockerDied = startDocker(origin, host, port, dclient)
		}
	}
}
