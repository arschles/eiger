package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/arschles/eiger/lib/util"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"io"
	"bytes"
)

func writeErr(err error, res http.ResponseWriter) {
	res.WriteHeader(http.StatusInternalServerError)
	body := fmt.Sprintf(`{"error":"%s"}`, err.Error())
	res.Write([]byte(body))
}

func agentsFunc(w http.ResponseWriter, r *http.Request) {

}

func readBytes(toRead int64, reader io.Reader) error {
	buf := bytes.NewBuffer([]byte{})
	n, err := io.Copy(buf, reader)
	if toRead != n {
		return fmt.Errorf("expected to read %d bytes, read %d", toRead, n)
	}
	return err
}

func heartbeatFunc(agents *Agents) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		for {
			err := readBytes(1, ws)
			if err != nil {
				util.LogWarnf("couldn't read from ws connection %s", ws)
				continue
			}
			agent := NewAgent(ws.Config().Origin, ws)
			//if the agent wasn't added we're already watching it
			if !agents.Add(*agent) {
				continue
			}
			//if the agent was added start watching it
			go func(ws *websocket.Conn, agent Agent) {
				for {
					time.Sleep(agents.hb)
					err := readBytes(1, ws)
					if err != nil {
						util.LogWarnf("removed agent %s", agent)
						agents.Remove(agent)
					}
				}
			}(ws, *agent)
		}
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
