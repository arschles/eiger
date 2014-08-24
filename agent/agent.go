package main

import (
	"fmt"
	"github.com/arschles/eiger/lib/cmd"
	"github.com/arschles/eiger/lib/util"
	"github.com/arschles/eiger/lib/ws"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"log"
	"time"
)

func agent(c *cli.Context) {
	dclient, err := docker.NewClient(c.String("dockerhost"))
	if err != nil {
		log.Fatal(err)
	}

	hbIntv := time.Millisecond * time.Duration(c.Int("heartbeat"))

	host := c.String("host")
	port := c.Int("port")
	socketUrl := fmt.Sprintf("ws://%s:%d/socket", host, port)
	heartUrl := fmt.Sprintf("ws://%s:%d/heart", host, port)
	origin := fmt.Sprintf("http://%s/", host)

	log.Printf("dialing %s", socketUrl)
	socketWs := ws.MustDial(socketUrl, origin)
	log.Printf("dialing %s", heartUrl)
	heartWs := ws.MustDial(heartUrl, origin)

	go heartbeater(heartWs, hbIntv)

	for {
		c, err := cmd.ReadCommand(socketWs)
		if err != nil {
			util.LogWarnf("reading a command: %s", err)
			continue
		}
		dispatch, ok := dispatchTable[c.Method]
		if !ok {
			go unknown(c, socketWs)
			continue
		}

		go dispatch(c, socketWs, dclient)
	}
}
