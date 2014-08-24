package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/arschles/eiger/lib/util"
	"log"
	"time"
)

func heartbeater(ws *websocket.Conn, hbIntv time.Duration) {
	numFailures := 0
	for {
		_, err := ws.Write([]byte{util.HeartbeatByte})
		if err != nil {
			numFailures++
			log.Printf("[WARN] heartbeat failure #%d: %s", numFailures, err)
		} else {
			numFailures = 0
		}
		time.Sleep(hbIntv)
	}
}
