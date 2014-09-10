package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/fsouza/go-dockerclient"
)

func logsLoop(ws *websocket.Conn, dockerClient *docker.Client, diedCh chan<- error) {

}
