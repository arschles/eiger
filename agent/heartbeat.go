package agent

import (
  "time"
  "encoding/json"
  "code.google.com/p/go.net/websocket"
)

type HeartbeatMessage struct {
  Hostname string `json:"hostname"`
  SendTime Time `json:"time"`
}

func heartbeat(wsConn *websocket.Conn, interval time.Duration, diedCh chan<- error) {
  hostname, err := os.Hostname()
  if err != nil {
    diedCh <- err
    return
  }

  for {
    msg := HeartbeatMessage{hostname, time.Now()}
    bytes, err := json.Marshal(msg)
    //TODO: backoff or fail if the heartbeat loop keeps erroring
    if err != nil {
      util.LogWarnf("(error heartbeating) %s", err)
    }

    sendStr := fmt.Sprintf("%d\n%s", len(bytes), bytes)
    sendBytes := []byte(sendStr)
    (numWritten, err) := wsConn.Write(sendBytes)
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
