package fig

/*
import (
    "strings"
    "github.com/gonuts/yaml"
    "github.com/fsouza/go-dockerclient"
    "fmt"
)

var emptyReader = strings.NewReader("")

type FigService struct {
    Image string `yaml:"image"`
    Cmd []string `yaml:"command"`
    Links []string `yaml:"links"`
    Ports []string `yaml:"ports"`
    Expose []string `yaml:"expose"`
    Env []string `yaml:"environment"`
    WorkingDir string `yaml:"working_dir"`
    Entrypoint []string `yaml:"entrypoint"`
}

type Fig struct {
    Services map[string]FigService
}

//ReadFig reads a fig yaml file (see http://www.fig.sh/yml.html for details)
//from figFile and creates a Fig struct from it. fig services without an "image"
//entry are omitted from the services that are returned in the Fig struct.
//the following elements are ignored by this function:
//{build, volumes, volumes_from, net, dns, user, hostname, domainname, mem_limit, privileged}
func Read(figFile []byte) (Fig, error) {
    var m map[string]FigService
    err := yaml.Unmarshal(figFile, &m)
    if err != nil {
        return Fig{}, err
    }
    newMap := map[string]FigService{}
    for name, svc := range m {
        if len(svc.Image) > 0 {
            newMap[name] = svc
        }
    }

    return Fig{newMap}, nil
}

func createFigContainer(name string, f FigService, d *docker.Client) (*docker.Container, error) {
    exposedPorts := map[docker.Port]struct{}{}
    for _, port := range f.Expose {
        //exposedPorts[port] = struct{}{}
    }
    conf := docker.Config {
        Image: f.Image,
        Env: f.Env,
        Entrypoint: f.Entrypoint,
        WorkingDir: f.WorkingDir,
        Cmd: f.Cmd,
        ExposedPorts: exposedPorts,
        //TODO: Links,Ports
    }
    cOpts := docker.CreateContainerOptions {
        Name: name,
        Config: &conf,
    }
    return d.CreateContainer(cOpts)
}

func Run(baseName string, f Fig, d *docker.Client) ([]*docker.Container, error) {
    names := map[string][]int{}
    containers := []*docker.Container{}
    for name, svc := range f.Services {
        ints, ok := names[name]
        newInt := 1
        if ok && len(ints) == 0 {
            newInt = 1
        } else if ok && len(ints) > 0 {
            lastInt := ints[len(ints)-1]
            newInt = lastInt + 1
        } else {
            newInt = 1
        }
        newInts := append(ints, newInt)
        names[name] = newInts
        //TODO: do this concurrently
        container, err := createFigContainer(fmt.Sprintf("%s%d", baseName, newInt), svc, d)
        if err != nil {
            return []*docker.Container{}, err
        }
        containers = append(containers, container)
    }
    return containers, nil
}
*/
