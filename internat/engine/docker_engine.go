package engine

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"

	"bigs-ci/lib/logger"
)

type DockerEngine struct {
	client   *client.Client
	hostPort string
	Ctx      context.Context
}

var dockerClient *client.Client

type ContainerConfig struct {
	ContainerName string
	ImageName     string
	WorkingDir    string
	Binds         []string
	VolumesFrom   []string
	PortBinds     []PortNat
	Cmd           string
	ExtraHosts    []string
	ExposedPorts  nat.PortSet
	PortBindings  nat.PortMap
	Restart       string
	Env           []string
}

type PortNat struct {
	Protocol   string `json:"protocol"`
	SourcePort string `json:"source_port"`
	TargetAddr string `json:"target_addr"`
	TargetPort string `json:"target_port"`
}

func New(ctx context.Context, hostPort string) *DockerEngine {
	cli := new(DockerEngine)
	if hostPort == "" {
		cli.hostPort = "tcp://192.168.3.21:2375"
	} else {
		cli.hostPort = hostPort
	}
	if cli.Ctx == nil {
		cli.Ctx = context.Background()
	} else {
		cli.Ctx = ctx
	}

	if dockerClient == nil {
		c, err := client.NewClientWithOpts(client.WithHost(cli.hostPort), client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		dockerClient = c
	}

	cli.client = dockerClient
	return cli
}

func (e *DockerEngine) Create(config *ContainerConfig, detach bool) (string, error) {
	containers, err := e.List(config.ContainerName, nil)
	if err != nil {
		logger.Error("queryContainerError", logger.Err(err))
		return "", err
	}
	if len(containers) != 0 {
		err := e.Remove(containers[0])
		if err != nil {
			logger.Error("removeContainerError", logger.Err(err))
			return "", err
		}
	}
	resp, err := e.client.ContainerCreate(e.Ctx, toConfig(config), toHostConfig(config), nil, nil, config.ContainerName)
	if err != nil {
		return "", err
	}

	if err := e.client.ContainerStart(e.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	if !detach {
		statusCh, errCh := e.client.ContainerWait(e.Ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				return "", err
			}
		case out := <-statusCh:
			if out.Error != nil {
				return "", errors.New(out.Error.Message)
			}
			if out.StatusCode != 0 {
				return "", errors.New("exec shell fail")
			}
		}

	}
	return resp.ID, nil
}

func (e *DockerEngine) Stop(ID string) error {
	return e.client.ContainerStop(e.Ctx, ID, nil)
}

func (e *DockerEngine) Restart(ID string) error {
	return e.client.ContainerRestart(e.Ctx, ID, nil)
}

func (e *DockerEngine) Remove(c types.Container) error {
	if c.State == "running" || c.State == "restarting" {
		err := e.Stop(c.ID)
		if err != nil {
			logger.Error("stopContainerError", logger.Err(err), logger.Any("container", c))
			return errors.New("stopContainerError")
		}
	}
	return e.client.ContainerRemove(e.Ctx, c.ID, types.ContainerRemoveOptions{})
}

func (e *DockerEngine) RemoveByID(ID string) error {
	return e.client.ContainerRemove(e.Ctx, ID, types.ContainerRemoveOptions{})
}

func (e *DockerEngine) Logs(containerName string) {
	out, err := e.client.ContainerLogs(e.Ctx, containerName, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
	if err != nil {
		logger.Error("printLogError", logger.Err(err))
		return
	}
	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		logger.Error("printLogToStdoutError", logger.Err(err))
		return
	}

}

func (e *DockerEngine) List(containerName string, status []string) ([]types.Container, error) {
	var args []filters.KeyValuePair
	if containerName != "" {
		args = []filters.KeyValuePair{filters.Arg("name", containerName)}

	}
	if len(status) == 0 {
		status = GetContainerStatus()
	}
	for _, m := range status {
		args = append(args, filters.Arg("status", m))
	}
	query := filters.NewArgs(args...)
	return e.client.ContainerList(e.Ctx, types.ContainerListOptions{Filters: query})
}

func (e *DockerEngine) Cleaner(containerName string) error {
	containers, err := e.List(containerName, nil)
	if err != nil {
		logger.Error("queryContainerError", logger.Err(err))
		return err
	}
	if len(containers) != 0 {
		err := e.Remove(containers[0])
		if err != nil {
			logger.Error("removeContainerError", logger.Err(err))
			return err
		}
	}
	return nil
}

func toConfig(config *ContainerConfig) *container.Config {

	c := &container.Config{
		Image:      config.ImageName,
		WorkingDir: config.WorkingDir,
		Env:        config.Env,
	}

	if len(config.PortBinds) != 0 && config.PortBinds[0].TargetPort != "" {
		portMap := make(map[nat.Port]struct{})
		for _, m := range config.PortBinds {
			portMap[nat.Port(m.SourcePort+"/"+m.Protocol)] = struct{}{}
		}
		c.ExposedPorts = portMap
	}

	if config.Cmd != "" {
		c.Cmd = toCmd(config.Cmd)
	}
	return c
}

func toHostConfig(config *ContainerConfig) *container.HostConfig {

	hc := &container.HostConfig{
		Binds:       config.Binds,
		VolumesFrom: config.VolumesFrom,
		ExtraHosts:  config.ExtraHosts,
	}

	if config.Restart == "always" {
		hc.RestartPolicy = container.RestartPolicy{
			Name: config.Restart,
		}
	}

	if len(config.PortBinds) != 0 && config.PortBinds[0].TargetPort != "" {
		portMap := make(map[nat.Port][]nat.PortBinding)
		for _, m := range config.PortBinds {
			portMap[nat.Port(fmt.Sprintf("%s/%s", m.SourcePort, m.Protocol))] = []nat.PortBinding{{HostIP: m.TargetAddr, HostPort: m.TargetPort}}
		}
		hc.PortBindings = portMap
	}

	return hc
}

func GetContainerStatus() []string {
	return []string{"created", "restarting", "running", "removing", "paused", "exited", "dead"}
}

func toCmd(command string) []string {
	return []string{"bash", "-c", command}

}
