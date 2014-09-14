package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/arschles/eiger/lib/messages"
	"github.com/arschles/eiger/lib/util"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

//the modulo value for printing heartbeat notifications
//TODO: make this configurable
const HBMOD = 10

//the multiplier for grace period between heartbeats
//TODO: make this configurable
const HBGRACEMULTIPLIER = 4

type heartbeatHandler struct {
	lookup     *AgentLookup
	hbInterval time.Duration
}

func (h *heartbeatHandler) serve(wsConn *websocket.Conn) {
	lastRecv := time.Now()
	iterNum := 0

	for {
		hbMsg := messages.Heartbeat{}
		err := websocket.JSON.Receive(wsConn, &hbMsg)
		if err != nil {
			util.LogWarnf("(parsing heartbeat message) %s", err)
			break
		}
		newAgent := NewAgent(hbMsg.Hostname, wsConn)
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
		lastRecv = time.Now()
	}
}

type dockerEventsHandler struct {
}

func (d *dockerEventsHandler) serve(wsConn *websocket.Conn) {
	for {
		time.Sleep(1 * time.Hour)
		//websocket.JSON.Receive(wsConn, dockerEvent)
	}
}

type rpcHandler struct {
}

func (r *rpcHandler) serve(ws *websocket.Conn) {
	for {
		time.Sleep(1 * time.Hour)
		//websocket.JSON.Receive(wsConn, rpcMethod)
	}
}

func service(c *cli.Context) {

	hbInterval := time.Duration(c.Int("heartbeat")) * time.Millisecond
	lookup := NewAgentLookup(&[]Agent{})

	hbHandler := heartbeatHandler{lookup, hbInterval}
	dockerEvtsHandler := dockerEventsHandler{}
	rpcHandler := rpcHandler{}

	router := mux.NewRouter()
	//REST verbs
	//r.HandleFunc("/agents", agentsFunc).Methods("GET")

	router.Handle("/heartbeat", websocket.Handler(hbHandler.serve))
	router.Handle("/docker", websocket.Handler(dockerEvtsHandler.serve))
	router.Handle("/rpc", websocket.Handler(rpcHandler.serve))

	ip := c.String("ip")
	port := c.Int("port")
	serveStr := fmt.Sprintf("%s:%d", ip, port)

	log.Printf("eiger-service listening on %s", serveStr)
	log.Fatal(http.ListenAndServe(serveStr, router))
}
