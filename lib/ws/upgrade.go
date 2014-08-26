package ws

import (
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/arschles/eiger/lib/util"
)

var Upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func Upgrade(f func(*websocket.Conn)) func(http.ResponseWriter, *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        conn, err := Upgrader.Upgrade(res, req, nil)
        if err != nil {
            util.LogWarnf("%s", err)
            return
        }
        f(conn)
    }
}
