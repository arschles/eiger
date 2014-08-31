package heartbeat

import (
    "io"
    "encoding/json"
    "time"
    "fmt"
    "github.com/arschles/eiger/lib/util"
    "strconv"
)

//HeartbeatMessage is an encoding.BinaryMarshaler that generates heartbeat
//messages to pass from agent to service to indicate the agent is still alive
type Message struct {
  Hostname string `json:"hostname"`
  SendTime time.Time `json:"time"`
}

func (h *Message) MarshalBinary() ([]byte, error) {
    bytes, err := json.Marshal(h)
    if err != nil {
        util.LogWarnf("(error heartbeating) %s", err)
        return []byte{}, err
    }
    sendStr := fmt.Sprintf("%d\n%s", len(bytes), err)
    return []byte(sendStr), nil
}

func parseLen(reader io.Reader) (int64, error) {
    //read until '\n' byte
    bytes := []byte{}
    for {
        b := make([]byte, 1)
        n, err := reader.Read(b)

        if err != nil {
            return 0, err
        } else if n <= 0 {
            return 0, fmt.Errorf("parsed length was %d", n)
        }

        if (n > 0 && b[0] == '\n') || n <= 0 {
            break
        }

        bytes = append(bytes, b[0])
    }

    return strconv.ParseInt(string(bytes), 10, 64)
}

//DecodeHeartbeatMessage reads the (simple) wire protocol generated by
//HeartbeatMessage.MarshalBinary into a HeartbeatMessage
func DecodeMessage(reader io.Reader) (*Message, error) {
    numBytes, err := parseLen(reader)
    if err != nil {
        return nil, err
    }
    bytes := make([]byte, numBytes)
    n, err := reader.Read(bytes)
    if err != nil {
        return nil, err
    }
    if int64(n) != numBytes {
        return nil, fmt.Errorf("expected to read %d bytes, but read %d", numBytes, n)
    }
    msg := new(Message)
    err = json.Unmarshal(bytes, msg)
    if err != nil {
        return nil, err
    }
    return msg, nil
}
