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

type socketHandler struct {
    lookup *AgentLookup
    hbLoop *HeartbeatLoop
}


func (h *socketHandler) serve(wsConn *websocket.Conn) {
    for {
        //TODO: have the heartbeat loop communicate back when the agent is dead
        hbMsg, err := heartbeat.DecodeMessage(wsConn)
        if err != nil {
            util.LogWarnf("(parsing heartbeat message) %s", err)
            return
        }
        newAgent := NewAgent(hbMsg.Hostname, wsConn)
        log.Printf("got agent %s", *newAgent)
        agent := h.lookup.GetOrAdd(*newAgent)
        h.hbLoop.Notify(*agent)
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

    socketHandler := socketHandler{set, hbLoop}

    router := mux.NewRouter()
    //REST verbs
    //r.HandleFunc("/agents", agentsFunc).Methods("GET")

    //Socket verb
    router.Handle("/socket", websocket.Handler(socketHandler.serve))

    //listen on websocket
    log.Fatal(http.ListenAndServe(serveStr, router))
}
