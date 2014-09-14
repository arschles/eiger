package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/arschles/eiger/lib/messages"
	"github.com/arschles/eiger/lib/util"
	"log"
	"os"
	"time"
)

//this is the modulo for heartbeat messages.
//TODO: make this configurable
const HBMOD = 10

//the number of consecutive heartbeat failures before dying
//TODO: make this configurable
const HBFAILTHRESH = 5

func heartbeatLoop(wsConn *websocket.Conn, interval time.Duration, diedCh chan<- error) {
	hostname, err := os.Hostname()
	if err != nil {
		diedCh <- err
		return
	}

	hbNum := 0
	lastFailed := false
	numFails := 0
	for {
		msg := messages.Heartbeat{
			Hostname: hostname,
			SendTime: time.Now(),
		}

		if hbNum%HBMOD == 0 {
			log.Printf("sending heartbeat message %s (%d)", msg, hbNum)
		}
		hbNum++

		err := websocket.JSON.Send(wsConn, msg)

		if err != nil {
			util.LogWarnf("(error heartbeating) %s", err)
			if lastFailed && numFails >= HBFAILTHRESH {
				break
			}
			numFails++
			lastFailed = true
		} else {
			lastFailed = false
			numFails = 0
		}

		time.Sleep(interval)
	}
	diedCh <- fmt.Errorf("heartbeat loop stopped")
}
