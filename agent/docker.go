package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/fsouza/go-dockerclient"
	"github.com/arschles/eiger/lib/messages"
	"github.com/arschles/eiger/lib/util"
)

func dockerLoop(ws *websocket.Conn, dockerClient *docker.Client, diedCh chan<- error) {
	evtChan := make(chan *docker.APIEvents)
	dockerClient.AddEventListener(evtChan)
	for evts := range evtChan {
		err := websocket.JSON.Send(ws, messages.DockerEvents(*evts))
		if err != nil {
			util.LogWarnf("(docker events) %s", err)
		}
	}
}
