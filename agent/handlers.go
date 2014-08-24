package main

import (
	"encoding/json"
	"fmt"
	"github.com/arschles/eiger/lib/cmd"
	"github.com/fsouza/go-dockerclient"
	"io"
)

type Handler func(*cmd.Command, io.Writer, *docker.Client)

var dispatchTable = map[string]Handler{
	cmd.DockerListContainersMethod:  listContainers,
	cmd.DockerVersionMethod:         version,
	cmd.DockerCreateContainerMethod: createContainer,
}

func unknown(c *cmd.Command, writer io.Writer) {
	c.WriteError(fmt.Errorf(cmd.UnknownMethod), writer)
}

func version(c *cmd.Command, writer io.Writer, d *docker.Client) {
	env, err := d.Version()
	if err != nil {
		c.WriteError(err, writer)
		return
	}
	c.WriteResult(env, nil, writer)
}

func listContainers(c *cmd.Command, writer io.Writer, d *docker.Client) {
	opts_string, ok := c.Params["opts_json"]
	if !ok {
		c.WriteError(fmt.Errorf("no opts_json field found"), writer)
		return
	}

	opts := docker.ListContainersOptions{}
	err := json.Unmarshal([]byte(opts_string), &opts)
	if err != nil {
		c.WriteError(err, writer)
		return
	}

	containers, err := d.ListContainers(opts)
	if err != nil {
		c.WriteError(err, writer)
		return
	}
	c.WriteResult(containers, nil, writer)
}

func createContainer(c *cmd.Command, writer io.Writer, d *docker.Client) {
	name, ok := c.Params["name"]
	if !ok {
		c.WriteError(fmt.Errorf("no name field found"), writer)
		return
	}
	opts_string, ok := c.Params["config_json"]
	if !ok {
		c.WriteError(fmt.Errorf("no config_json field found"), writer)
		return
	}
	conf := docker.Config{}
	err := json.Unmarshal([]byte(opts_string), &conf)
	if err != nil {
		c.WriteError(err, writer)
		return
	}
	opts := docker.CreateContainerOptions{
		Name:   name,
		Config: &conf,
	}

	container, err := d.CreateContainer(opts)
	if err != nil {
		c.WriteError(err, writer)
		return
	}
	c.WriteResult(container, nil, writer)
}
