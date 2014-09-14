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
			Value:  2000,
			Usage:  "the longest allowable latency between heartbeats, in milliseconds",
			EnvVar: "EIGER_HEARTBEAT",
		},
		cli.StringSliceFlag{
			Name:   "publishtypes,pts",
			Value:  &cli.StringSlice{"log"},
			Usage:  "the methods by which the eiger service should publish incoming data",
			EnvVar: "EIGER_PUBLISH_TYPES",
		},
	}
	app.Version = "0.1.0"
	app.Action = service
	app.Run(os.Args)
}
