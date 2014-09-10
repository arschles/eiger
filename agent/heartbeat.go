package main

import (
  "time"
  "code.google.com/p/go.net/websocket"
  "os"
  "github.com/arschles/eiger/lib/util"
  "github.com/arschles/eiger/lib/heartbeat"
  "fmt"
  "log"
)

func heartbeatLoop(wsConn *websocket.Conn, interval time.Duration, diedCh chan<- error) {
  hostname, err := os.Hostname()
  if err != nil {
    diedCh <- err
    return
  }

  for {
    msg := heartbeat.Message{hostname, time.Now()}
    log.Printf("sending heartbeat message %s", msg)
    bytes, err := msg.MarshalBinary()
    //TODO: backoff or fail if the heartbeat loop keeps erroring
    if err != nil {
      util.LogWarnf("(error heartbeating) %s", err)
    }

    sendStr := fmt.Sprintf("%d\n%s", len(bytes), bytes)
    sendBytes := []byte(sendStr)
    numWritten, err := wsConn.Write(sendBytes)
    if numWritten != len(sendBytes) {
      util.LogWarnf("(error heartbeating) wrote %d bytes, expected %d", numWritten, len(sendBytes))
    }
    if err != nil {
      util.LogWarnf("(error heartbeating) %s", err)
    }
    time.Sleep(interval)
  }
  diedCh <- fmt.Errorf("heartbeat loop stopped")
}
