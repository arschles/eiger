package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"os"
	"time"
)

func run(c *cli.Context) {
	ip := c.String("ip")
	port := c.Int("port")
	serveStr := fmt.Sprintf("%s:%d", ip, port)
	log.Printf("eiger-service listening on %s", serveStr)

	dur := time.Duration(c.Int("heartbeat")) * time.Millisecond
	router := router(dur)
	log.Fatal(http.ListenAndServe(serveStr, router))
}

func main() {
	app := cli.NewApp()
	app.Name = "eiger-service"
	app.Usage = "connection service between eiger-agent and the eiger backend"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ip,i",
			Value:  "127.0.0.1",
			Usage:  "the IP address to listen on",
			EnvVar: "EIGER_HOST",
		},
		cli.IntFlag{
			Name:   "port,p",
			Value:  4492,
			Usage:  "the port to listen on",
			EnvVar: "EIGER_PORT",
		},
		cli.IntFlag{
			Name:   "heartbeat,b",
			Value:  1000,
			Usage:  "the longest allowable latency between heartbeats",
			EnvVar: "EIGER_HEARTBEAT",
		},
	}
	app.Version = "0.1.0"
	app.Action = run
	app.Run(os.Args)
}
