package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/arschles/eiger/lib/heartbeat"
	"github.com/arschles/eiger/lib/util"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type socketHandler struct {
	lookup     *AgentLookup
	hbInterval time.Duration
}

//the modulo value for printing heartbeat notifications
//TODO: make this configurable
const HBMOD = 10

func (h *socketHandler) serve(wsConn *websocket.Conn) {
	watcherRunning := false
	tickerCh := make(chan interface{})
	removedCh := make(chan interface{})
	iterNum := 0
	for {
		//TODO: have the heartbeat loop communicate back when the agent is dead
		hbMsg, err := heartbeat.DecodeMessage(wsConn)
		if err != nil {
			util.LogWarnf("(parsing heartbeat message) %s", err)
			break
		}

		newAgent := NewAgent(hbMsg.Hostname, wsConn)
		agent := h.lookup.GetOrAdd(*newAgent)
		if iterNum%HBMOD == 0 {
			log.Printf("got agent %s", *agent)
		}
		iterNum++

		if !watcherRunning {
			log.Printf("starting watch loop for agent %s", *agent)
			go agentWatcher(*agent, h.hbInterval, tickerCh, removedCh)
		}

		tickerCh <- struct{}{}
	}
}

func service(c *cli.Context) {
	ip := c.String("ip")
	port := c.Int("port")
	serveStr := fmt.Sprintf("%s:%d", ip, port)
	log.Printf("eiger-service listening on %s", serveStr)

	hbInterval := time.Duration(c.Int("heartbeat")) * time.Millisecond
	set := NewAgentLookup(&[]Agent{})

	socketHandler := socketHandler{set, hbInterval}

	router := mux.NewRouter()
	//REST verbs
	//r.HandleFunc("/agents", agentsFunc).Methods("GET")

	//Socket verb
	router.Handle("/socket", websocket.Handler(socketHandler.serve))

	//listen on websocket
	log.Fatal(http.ListenAndServe(serveStr, router))
}
