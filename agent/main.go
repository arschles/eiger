package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "eiger-agent"
	app.Usage = "phone home to the eiger service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host,o",
			Value:  "127.0.0.1",
			Usage:  "the eiger host to connect to",
			EnvVar: "EIGER_HOST",
		},
		cli.IntFlag{
			Name:   "port,p",
			Value:  4492,
			Usage:  "the port for to connect to the host on",
			EnvVar: "EIGER_PORT",
		},
		cli.IntFlag{
			Name:   "heartbeat,b",
			Value:  1000,
			Usage:  "the time between heartbeats, in milliseconds",
			EnvVar: "EIGER_HEARTBEAT",
		},
		cli.StringFlag{
			Name:   "dockerhost,d",
			Value:  "/var/run/docker.sock",
			Usage:  "the docker host to connect to",
			EnvVar: "DOCKER_HOST",
		},
	}
	app.Version = "0.1.0"
	app.Action = agent
	app.Run(os.Args)
}
