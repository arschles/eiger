package main

import (
    "fmt"
    "github.com/codegangsta/cli"
    "log"
    "net/http"
    "time"
    "code.google.com/p/go.net/websocket"
    "github.com/gorilla/mux"
    "github.com/arschles/eiger/lib/heartbeat"
    "github.com/arschles/eiger/lib/util"
)

type heartbeatHandler struct {
    lookup *AgentLookup
    hbLoop *HeartbeatLoop
}


func (h *heartbeatHandler) serve(wsConn *websocket.Conn) {
    for {
        //TODO: have the heartbeat loop communicate back when the agent is dead
        hbMsg, err := heartbeat.DecodeMessage(wsConn)
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
        //websocket.JSON.Receive(wsConn, dockerEvent)
    }
}

func service(c *cli.Context) {
    ip := c.String("ip")
    port := c.Int("port")
    serveStr := fmt.Sprintf("%s:%d", ip, port)
    log.Printf("eiger-service listening on %s", serveStr)

    hbDur := time.Duration(c.Int("heartbeat")) * time.Millisecond
    set := NewAgentLookup(&[]Agent{})
    hbLoop := NewHeartbeatLoop(set, hbDur)

    hbHandler := heartbeatHandler{set, hbLoop}
    dockerEvtsHandler := dockerEventsHandler{}

    router := mux.NewRouter()
    //REST verbs
    //r.HandleFunc("/agents", agentsFunc).Methods("GET")

    router.Handle("/heartbeat", websocket.Handler(hbHandler.serve))
    router.Handle("/docker_events", websocket.Handler(dockerEvtsHandler.serve))

    //listen on websocket
    log.Fatal(http.ListenAndServe(serveStr, router))
}
