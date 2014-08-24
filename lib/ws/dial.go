package ws

import (
    "code.google.com/p/go.net/websocket"
    "log"
)

//MustDial dials url at origin and returns the resulting connection.
//If there's an error, calls log.Fatal
func MustDial(url, origin string) *websocket.Conn {
    ws, err := websocket.Dial(url, "", origin)
    if err != nil {
        log.Fatal(err)
    }
    return ws
}
