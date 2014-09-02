package main

import (
    "fmt"
    "github.com/codegangsta/cli"
    "log"
    "net/http"
    "time"
    "code.google.com/p/go.net/websocket"
    "net/rpc"
    "net/rpc/jsonrpc"
    "github.com/gorilla/mux"
)

type socketHandler struct {
    rpcServer *rpc.Server
}

func (h *socketHandler) serve(wsConn *websocket.Conn) {
    serverCodec := jsonrpc.NewServerCodec(wsConn)
    h.rpcServer.ServeCodec(serverCodec)
}

type heartbeatHandler struct {
    agentSet AgentSet
}

func (h *heartbeatHandler) serve(wsConn *websocket.Conn) {
    //TODO: add to agent set
}

func service(c *cli.Context) {
    ip := c.String("ip")
    port := c.Int("port")
    serveStr := fmt.Sprintf("%s:%d", ip, port)
    log.Printf("eiger-service listening on %s", serveStr)

    heartbeat := time.Duration(c.Int("heartbeat")) * time.Millisecond

    set := NewAgentSet(&[]Agent{})

    rpcServer := rpc.NewServer()
    rpcServer.Register(NewHandlers(set, heartbeat))
    sockHandler := socketHandler{rpcServer}

    router := mux.NewRouter()
    //REST verbs
    //r.HandleFunc("/agents", agentsFunc).Methods("GET")

    //Socket verb
    router.Handle("/socket", websocket.Handler(sockHandler.serve))

    //listen on websocket
    log.Fatal(http.ListenAndServe(serveStr, router))
}
