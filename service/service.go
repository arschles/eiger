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

type SocketHandler struct {
    rpcServer *rpc.Server
    connLookup *ConnLookup
}

func NewSocketHandler(rpcServer *rpc.Server, connLookup *ConnLookup) *SocketHandler) {
    return &SocketHandler{
        rpcServer: rpcServer,
        connLookup: connLookup,
    }
}

func (h *socketHandler) serve(wsConn *websocket.Conn) {
    connLookup.AddUnmatched(wsConn)
    serverCodec := jsonrpc.NewServerCodec(wsConn)
    h.rpcServer.ServeCodec(serverCodec)
}

func service(c *cli.Context) {
    ip := c.String("ip")
    port := c.Int("port")
    serveStr := fmt.Sprintf("%s:%d", ip, port)
    log.Printf("eiger-service listening on %s", serveStr)

    heartbeat := time.Duration(c.Int("heartbeat")) * time.Millisecond

    set := NewAgentSet(&[]Agent{})
    connLookup := NewConnLookup()
    handlers := NewHandlers(set, connLookup, heartbeat)

    rpcServer := rpc.NewServer()
    rpcServer.Register(handlers)
    sockHandler := NewSocketHandler(rpcServer, connLookup)

    router := mux.NewRouter()
    //REST verbs
    //r.HandleFunc("/agents", agentsFunc).Methods("GET")

    //Socket verb
    router.Handle("/socket", websocket.Handler(sockHandler.serve))

    //listen on websocket
    log.Fatal(http.ListenAndServe(serveStr, router))
}
