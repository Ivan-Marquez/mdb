package container

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// Container represents a Docker container and its configuration
type Container struct {
	RemoteURL       string                // image repository URL
	Attachments     []string              // files you want to copy to the container
	HostConfig      *container.HostConfig // host configuration
	ContainerConfig *container.Config     // container configuration
	ExecConfig      *types.ExecConfig     // runtime configuration to execute once the container has started
}

// ID represents the running container Id
type ID string
