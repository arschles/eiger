package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/arschles/eiger/lib/messages"
	"github.com/arschles/eiger/lib/pubsub"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

func publishAll(payload *pubsub.Payload, publishers []pubsub.Publisher) {
	for _, publisher := range publishers {
		publisher.Publish(payload)
	}
}

type heartbeatHandler struct {
	lookup     *AgentLookup
	hbInterval time.Duration
	publishers []pubsub.Publisher
}

func (h *heartbeatHandler) serve(wsConn *websocket.Conn) {
	lastRecv := time.Now()
	iterNum := 0
	numLate := 0
	for {
		hbMsg := messages.Heartbeat{}
		err := websocket.JSON.Receive(wsConn, &hbMsg)
		if err != nil {
			payload := pubsub.Payload{pubsub.HeartbeatErrorTopic, err}
			publishAll(&payload, h.publishers)
			break
		}
		newAgent := NewAgent(hbMsg.Hostname, wsConn)
		agent := h.lookup.GetOrAdd(*newAgent)
		if time.Since(lastRecv) > (h.hbInterval * HeartbeatGraceMultiplier) {
			payload := pubsub.Payload{pubsub.LateHeartbeatTopic, *agent}
			publishAll(&payload, h.publishers)
			h.lookup.Remove(*agent)
			if numLate > MaxNumLateHeartbeats {
				break
			} else {
				numLate++
				iterNum++
				continue
			}
		}

		if iterNum%HeartbeatModulus == 0 {
			payload := pubsub.Payload{pubsub.HeartbeatTopic, map[string]string {
				"heartbeat_num": fmt.Sprintf("%d", iterNum),
				"agent": fmt.Sprintf("%s", *agent),
			}}
			publishAll(&payload, h.publishers)
		}

		iterNum++

		lastRecv = time.Now()
	}
}

type dockerEventsHandler struct {
	publishers []pubsub.Publisher
}

func (d *dockerEventsHandler) serve(wsConn *websocket.Conn) {
	for {
		evts := messages.DockerEvents{}
		err := websocket.JSON.Receive(wsConn, &evts)
		payload := pubsub.Payload{pubsub.DockerEventsTopic, evts}
		if err != nil {
			payload = pubsub.Payload{pubsub.DockerEventsErrorTopic, err}
		}
		publishAll(&payload, d.publishers)
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

func parsePublishers(slice []string) []pubsub.Publisher {
	pslice := []pubsub.Publisher{}
	for _, str := range slice {
		switch str {
		case "log":
			pslice = append(pslice, pubsub.LoggingPublisher{})
		}
	}
	return pslice
}

func service(c *cli.Context) {

	publishers := parsePublishers(c.StringSlice("publishtypes"))

	hbInterval := time.Duration(c.Int("heartbeat")) * time.Millisecond
	lookup := NewAgentLookup(&[]Agent{})

	hbHandler := heartbeatHandler{lookup, hbInterval, publishers}
	dockerEvtsHandler := dockerEventsHandler{publishers}
	rpcHandler := rpcHandler{}

	router := mux.NewRouter()

	router.Handle("/heartbeat", websocket.Handler(hbHandler.serve))
	router.Handle("/docker", websocket.Handler(dockerEvtsHandler.serve))
	router.Handle("/rpc", websocket.Handler(rpcHandler.serve))

	ip := c.String("ip")
	port := c.Int("port")
	serveStr := fmt.Sprintf("%s:%d", ip, port)

	log.Printf("eiger-service listening on %s", serveStr)
	log.Fatal(http.ListenAndServe(serveStr, router))
}
