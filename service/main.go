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
			Name:   "service-host,sh",
			Value:  "127.0.0.1",
			Usage:  "the host for the service to listen on",
			EnvVar: "EIGER_SERVICE_HOST",
		},
		cli.IntFlag{
			Name:   "service-port,sp",
			Value:  4492,
			Usage:  "the port for the service to listen on",
			EnvVar: "EIGER_SERVICE_PORT",
		},
		cli.IntFlag{
			Name:   "service-heartbeat,b",
			Value:  2000,
			Usage:  "the longest duration (milliseconds) that the service will wait to hear back from an agent",
			EnvVar: "EIGER_SERVICE_HEARTBEAT",
		},
		cli.StringSliceFlag{
			Name:   "service-publish-types,spt",
			Value:  &cli.StringSlice{"log"},
			Usage:  "the methods by which the service should publish incoming data",
			EnvVar: "EIGER_SERVICE_PUBLISH_TYPES",
		},
		cli.StringFlag{
			Name: "api-host,ah",
			Value: "127.0.0.1",
			Usage: "the host for the api to listen on",
			EnvVar: "EIGER_API_HOST",
		},
		cli.IntFlag{
			Name:"api-port,ap",
			Value:4493,
			Usage:"the port for the api to listen on",
			EnvVar:"EIGER_API_PORT",
		},
		cli.StringSliceFlag{
			Name:"api-subscribe-types,ast",
			Value: &cli.StringSlice{"ptp-udp"},
			Usage: "the methods by which the api should listen for published data",
			EnvVar:"EIGER_API_SUBSCRIBE_TYPES",
		},
	}
	app.Version = "0.1.0"
	app.Action = func(c *cli.Context) {
		ch := make(chan interface{})

		go func() {
			service(c)
			ch <- struct{}{}
		}()
		go func() {
			api(c)
			ch <- struct{}{}
		}()
		<-ch
	}
	app.Run(os.Args)
}
