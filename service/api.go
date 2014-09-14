package main

import (
  "github.com/codegangsta/cli"
  "github.com/gorilla/mux"
  "net/http"
  "log"
  "fmt"
)

func api(c *cli.Context) {
  router := mux.NewRouter()
  host := c.String("api-host")
  port := c.Int("api-port")
  serveStr := fmt.Sprintf("%s:%d", host, port)
  log.Printf("eiger-api listening on %s", serveStr)
  log.Fatal(http.ListenAndServe(serveStr, router))
}
