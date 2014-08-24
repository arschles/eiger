package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/arschles/eiger/lib/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func writeErr(err error, res http.ResponseWriter) {
	res.WriteHeader(http.StatusInternalServerError)
	body := fmt.Sprintf(`{"error":"%s"}`, err.Error())
	res.Write([]byte(body))
}

func agentsFunc(w http.ResponseWriter, r *http.Request) {

}

func heartbeatFunc(agents *Agents) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		log.Printf("got heartbeat from %s", ws.Config().Origin)
		agent := NewAgent(ws.Config().Origin, ws)
		ch := make(chan Agent)
		agents.Add(*agent, ch)
		go func() {
			removed := <-ch
			util.LogWarnf("removed agent %s", removed.Origin)
		}()
	}
}

func socketFunc(ws *websocket.Conn) {

}

func router(hb time.Duration) *mux.Router {
	agents := NewAgents(&[]Agent{}, hb)

	r := mux.NewRouter()
	r.HandleFunc("/agents", agentsFunc).Methods("GET")
	r.Handle("/heart", websocket.Handler(heartbeatFunc(agents)))
	r.Handle("/socket", websocket.Handler(socketFunc))
	return r
}
