package cmd

import (
    "code.google.com/p/go-uuid/uuid"
    "github.com/fsouza/go-dockerclient"
    "encoding/json"
    "io"
    "fmt"
    "bytes"
)

const (
    DockerListContainersMethod = "d_list_containers"
    DockerVersionMethod = "d_version"
    DockerCreateContainerMethod = "d_create_container"

    UnknownMethod = "e_unknown_method"
)

//Command represents a (json) message the server sends to the agent to do
//some action on its behalf. Command is json-rpc compliant.
type Command struct {
    Method string `json:"method"`
    Params map[string]string `json:"args"`
    Id string `json:"id"`
}

//ReadCommand blocks until it reads >0 bytes from reader (using io.Copy),
//then attempts to deserialize those bytes into a Command. when it finds
//an error, returns nil, err. otherwise returns the deserialized command, nil
func ReadCommand(reader io.Reader) (*Command, error) {
    cmd := Command{}
    buf := bytes.NewBuffer([]byte{})
    for {
        n, err := io.Copy(buf, reader)
        if err != nil {
            return nil, err
        }
        if n > 0 {
            break
        }
    }
    err := json.Unmarshal(buf.Bytes(), &cmd)
    if err != nil {
        return nil, err
    }
    return &cmd, nil
}

//WriteResult writes the given result and error to the given writer, returning
//an encoding or send error if one occurred.
//the result will be json-rpc compliant, so exactly one of result or err must
//be nil. result must be marshalable with the encoding/json package.
func (c *Command) WriteResult(result interface{}, err error, writer io.Writer) error {
    m := map[string]interface{} {
        "result": result,
        "error": err,
        "id": c.Id,
    }
    bytes, err := json.Marshal(m)
    if err != nil {
        return err
    }
    n, err := writer.Write(bytes)
    if err != nil {
        return err
    }
    if n != len(bytes) {
        return fmt.Errorf("wrote %d bytes but expected to write %d", n, len(bytes))
    }
    return nil
}

//WriteError is a convenience method for c.WriteResult(nil, err, writer)
func (c *Command) WriteError(err error, writer io.Writer) error {
    return c.WriteResult(nil, err, writer)
}

//DockerListContainers returns a command to tell the agent to send back a list
//of containers
func DockerListContainers(opts *docker.ListContainersOptions) (*Command, error) {
    bytes, err := json.Marshal(opts)
    if err != nil {
        return nil, err
    }
    cmd := Command{
        Method: DockerListContainersMethod,
        Params: map[string]string{
            "opts_json": string(bytes),
        },
        Id: uuid.New(),
    }

    return &cmd, nil
}

//DockerListContainers returns a command to tell the agent to send back the
//docker version
func DockerVersion() *Command {
    return &Command{
        Method: DockerVersionMethod,
        Id: uuid.New(),
    }
}

//DockerCreateContainer returns a command to tell the agent to create a
//container and return the result
func DockerCreateContainer(name string, config *docker.Config) (*Command, error) {
    bytes, err := json.Marshal(config)
    if err != nil {
        return nil, err
    }

    cmd := Command{
        Method: DockerCreateContainerMethod,
        Params: map[string]string{
            "name": name,
            "config_json": string(bytes),
        },
        Id: uuid.New(),
    }
    return &cmd, nil
}
