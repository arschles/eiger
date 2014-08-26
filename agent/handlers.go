package main

import (
	"github.com/fsouza/go-dockerclient"
)

type Handlers struct {
	dclient *docker.Client
}

func NewHandlers(d *docker.Client) *Handlers {
	return &Handlers{d}
}

func (h *Handlers) Heartbeat(_ struct{}, reply *int) error {
	*reply = 0
	return nil
}

func (h *Handlers) ListContainers(opts docker.ListContainersOptions, reply *[]docker.APIContainers) error {
	containers, err := h.dclient.ListContainers(opts)
	if err != nil {
		return err
	}
	*reply = containers
	return nil
}

func (h *Handlers) Version(a struct{}, reply *docker.Env) error {
	env, err := h.dclient.Version()
	if err != nil {
		return err
	}
	*reply = *env
	return nil
}

func (h *Handlers) CreateContainer(opts docker.CreateContainerOptions, reply *docker.Container) error {
	container, err := h.dclient.CreateContainer(opts)
	if err != nil {
		return err
	}
	*reply = *container
	return nil
}
