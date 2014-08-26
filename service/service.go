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

func service(c *cli.Context) {
    ip := c.String("ip")
    port := c.Int("port")
    serveStr := fmt.Sprintf("%s:%d", ip, port)
    log.Printf("eiger-service listening on %s", serveStr)

    dur := time.Duration(c.Int("heartbeat")) * time.Millisecond

    set := NewAgentSet(&[]Agent{}, dur)

    rpcServer := rpc.NewServer()
    rpcServer.Register(NewHandlers(set))
    sockHandler := socketHandler{rpcServer}

    router := mux.NewRouter()
    //REST verbs
    //r.HandleFunc("/agents", agentsFunc).Methods("GET")

    //Socket verb
    router.Handle("/socket", websocket.Handler(sockHandler.serve))

    //listen on websocket
    log.Fatal(http.ListenAndServe(serveStr, router))
}

// func writeErr(err error, res http.ResponseWriter) {
// 	res.WriteHeader(http.StatusInternalServerError)
// 	body := fmt.Sprintf(`{"error":"%s"}`, err.Error())
// 	res.Write([]byte(body))
// }

// func heartbeatFunc(agents *Agents) func(*websocket.Conn) {
// 	return func(ws *websocket.Conn) {
// 		agent := NewAgent()
// 		if !agents.Add(agent) {
// 			return
// 		}
// 		ch := make(chan string)
// 		ws.SetPingHandler(func(s string) error {
// 			ch <- s
// 			return nil
// 		})
// 		go func() {
// 			for {
// 				select {
// 				case <-ch:
// 				case <-time.After(agents.hb):
// 					util.LogWarnf("didn't receive heartbeat in %s from agent %s. removing agent", agents.hb, agent)
// 					return
// 				}
// 			}
// 		}()
// 	}
// }
