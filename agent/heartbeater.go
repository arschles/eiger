package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"time"
	"github.com/arschles/eiger/lib/util"
)

func heartbeater(ws *websocket.Conn, hbIntv time.Duration) {
	numFailures := 0
	for {
		n, err := ws.Write([]byte{'h'})
		if err != nil {
			numFailures++
			util.LogWarnf("heartbeat failure #%d: %s", numFailures, err)
		} else {
			numFailures = 0
			log.Printf("sent %d bytes in heartbeat", n)
		}

		time.Sleep(hbIntv)
	}
}
