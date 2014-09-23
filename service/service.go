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
	"runtime"
	"time"
)

//the modulo value for printing heartbeat notifications
//TODO: make this configurable
const HeartbeatModulus = 10

//the multiplier for grace period between heartbeats
//TODO: make this configurable
const HeartbeatGraceMultiplier = 4

//the maximum number of late heartbeats allowed before shutting down the
//heartbeat connection.
//TODO: make this configurable
const MaxNumLateHeartbeats = 10

//the maximum number of errors receiving docker events allowed
//TODO: make this configurable
const MaxNumDockerEventErrors = 10

type heartbeatHandler struct {
	lookup     *AgentLookup
	hbInterval time.Duration
	streamChan chan<-interface{}
}

func (h *heartbeatHandler) serve(wsConn *websocket.Conn) {
	lastRecv := time.Now()
	iterNum := 0
	numLate := 0
	for {
		hbMsg := messages.Heartbeat{}
		err := websocket.JSON.Receive(wsConn, &hbMsg)
		if err != nil {
			h.streamChan <- err
			break
		}
		newAgent := NewAgent(hbMsg.Hostname, wsConn)
		agent := h.lookup.GetOrAdd(*newAgent)
		if time.Since(lastRecv) > (h.hbInterval * HeartbeatGraceMultiplier) {
			h.streamChan <- fmt.Sprintf("removed %s", *agent)
			h.lookup.Remove(*agent)
			lastRecv = time.Now()
			if numLate > MaxNumLateHeartbeats {
				break
			} else {
				numLate++
				iterNum++
				continue
			}
		}

		payload := map[string]string{
			"heartbeat_num": fmt.Sprintf("%d", iterNum),
			"agent":         fmt.Sprintf("%s", *agent),
		}
		h.streamChan <- payload

		iterNum++

		lastRecv = time.Now()
	}
}

type dockerEventsHandler struct {
	streamChan chan<-interface{}
}

func (d *dockerEventsHandler) serve(wsConn *websocket.Conn) {
	numErrs := 0
	for {
		evts := messages.DockerEvents{}
		err := websocket.JSON.Receive(wsConn, &evts)
		if err != nil {
			d.streamChan <- err
			if numErrs > MaxNumDockerEventErrors {
				break
			}
			numErrs++
			continue
		}
		d.streamChan <- evts
	}
}

type rpcHandler struct {
	streamChan chan<- interface{}
}

func (r *rpcHandler) serve(ws *websocket.Conn) {
	for {
		runtime.Gosched()
		//websocket.JSON.Receive(wsConn, rpcMethod)
	}
}

type streamHandler struct {
	b *util.Broadcaster
}

func (s *streamHandler) serve(ws *websocket.Conn) {
	ch := s.b.NewChan()
	for {
		evt := <-ch
		log.Printf("%s", evt)
		websocket.JSON.Send(ws, evt)
	}
}


func service(c *cli.Context) {

	hbInterval := time.Duration(c.Int("heartbeat")) * time.Millisecond
	lookup := NewAgentLookup(&[]Agent{})

	streamChan := make(chan interface{})
	broadcaster := util.NewBroadcaster(streamChan)

	hbHandler := heartbeatHandler{lookup, hbInterval, streamChan}
	dockerEvtsHandler := dockerEventsHandler{streamChan}
	rpcHandler := rpcHandler{streamChan}
	streamHandler := streamHandler{broadcaster}

	router := mux.NewRouter()

	router.Handle("/heartbeat", websocket.Handler(hbHandler.serve))
	router.Handle("/docker", websocket.Handler(dockerEvtsHandler.serve))
	router.Handle("/rpc", websocket.Handler(rpcHandler.serve))

	//this is stuff that the dashboard connects on
	router.PathPrefix("/stream").Handler(websocket.Handler(streamHandler.serve))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dashboard")))

	host := c.String("host")
	port := c.Int("port")
	serveStr := fmt.Sprintf("%s:%d", host, port)

	//start the logging goroutine
	go func() {
		ch := broadcaster.NewChan()
		n := 0
		for {
			value := <-ch
			switch t := value.(type) {
			case messages.DockerEvents:
				log.Printf("%s", t)
			default:
				log.Printf("%s", t)
			}
			n++
		}
	}()

	log.Printf("eiger-service listening on %s", serveStr)
	log.Fatal(http.ListenAndServe(serveStr, router))
}
