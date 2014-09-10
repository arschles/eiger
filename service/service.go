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

type heartbeatHandler struct {
	lookup *AgentLookup
	hbLoop *HeartbeatLoop
}

func (h *heartbeatHandler) serve(wsConn *websocket.Conn) {
	for {
		//TODO: have the heartbeat loop communicate back when the agent is dead
		hbMsg := messages.Heartbeat{}
		err := websocket.JSON.Receive(wsConn, &hbMsg)
		if err != nil {
			util.LogWarnf("(parsing heartbeat message) %s", err)
			return
		}
		newAgent := NewAgent(hbMsg.Hostname, wsConn)
		agent := h.lookup.GetOrAdd(*newAgent)
		h.hbLoop.Notify(*agent)
	}
}

type dockerEventsHandler struct {
}

func (d *dockerEventsHandler) serve(wsConn *websocket.Conn) {
	for {
		time.Sleep(1*time.Hour)
		//websocket.JSON.Receive(wsConn, dockerEvent)
	}
}

type rpcHandler struct {

}

func (r *rpcHandler) serve(ws *websocket.Conn) {
	for {
		time.Sleep(1*time.Hour)
		//websocket.JSON.Receive(wsConn, rpcMethod)
	}
}

func service(c *cli.Context) {
	hbDur := time.Duration(c.Int("heartbeat")) * time.Millisecond
	set := NewAgentLookup(&[]Agent{})
	hbLoop := NewHeartbeatLoop(set, hbDur)

	hbHandler := heartbeatHandler{set, hbLoop}
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
