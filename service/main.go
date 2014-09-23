package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "eiger-service"
	app.Usage = "connection service between eiger-agent and the eiger backend"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host,sh",
			Value:  "127.0.0.1",
			Usage:  "the host for the service to listen on",
			EnvVar: "EIGER_SERVICE_HOST",
		},
		cli.IntFlag{
			Name:   "port,sp",
			Value:  4492,
			Usage:  "the port for the service to listen on",
			EnvVar: "EIGER_SERVICE_PORT",
		},
		cli.IntFlag{
			Name:   "heartbeat,b",
			Value:  2000,
			Usage:  "the longest duration (milliseconds) that the service will wait to hear back from an agent",
			EnvVar: "EIGER_SERVICE_HEARTBEAT",
		},
	}
	app.Version = "0.1.0"
	app.Action = service
	app.Run(os.Args)
}
