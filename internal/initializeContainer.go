package internal

import (
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"

	mdb "github.com/ivan-marquez/mdb/pkg/container"
)

// InitializeMDBContainer initializes a Docker container
// with a MongoDB image
func InitializeMDBContainer(cmd, paths []string) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cc := &container.Config{
		WorkingDir:   dir,
		Image:        "mongo:latest",
		AttachStdout: true,
		AttachStderr: true,
	}

	hc := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: dir,
				Target: dir,
			},
		},
	}

	ec := &types.ExecConfig{
		Env:          []string{},
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	c := &mdb.Container{
		RemoteURL:       "docker.io/library/mongo",
		Attachments:     paths,
		ContainerConfig: cc,
		HostConfig:      hc,
		ExecConfig:      ec,
	}

	err = c.RunContainer()
	if err != nil {
		return err
	}

	return nil
}
