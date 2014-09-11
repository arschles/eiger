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

//the multiplier for grace period between heartbeats
//TODO: make this configurable
const HBGRACEMULTIPLIER = 4

func (h *socketHandler) serve(wsConn *websocket.Conn) {

	lastRecv := time.Now()
	iterNum := 0

	for {
		recvCh := make(chan heartbeat.Message)
		errCh := make(chan error)
		//TODO: have a way to kill this goroutine
		go func() {
			hbMsg, err := heartbeat.DecodeMessage(wsConn)
			if err != nil {
				errCh <- err
				return
			}
			recvCh <- *hbMsg
		}()

		select {
		case err := <-errCh:
			util.LogWarnf("(parsing heartbeat message) %s", err)
			break
		case recv := <-recvCh:
			newAgent := NewAgent(recv.Hostname, wsConn)
			agent := h.lookup.GetOrAdd(*newAgent)
			if iterNum%HBMOD == 0 {
				log.Printf("got agent heartbeat %s", *agent)
			}
			iterNum++
			if time.Since(lastRecv) > (h.hbInterval * HBGRACEMULTIPLIER) {
				util.LogWarnf("(late heartbeat) removing agent %s from alive set", *agent)
				h.lookup.Remove(*agent)
				break
			}
		}
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
