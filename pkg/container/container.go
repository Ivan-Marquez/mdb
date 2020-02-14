package container

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// contains tells whether a string slice contains specified string
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// imageExists checks if a Docker image exists within the host
func imageExists(ctx context.Context, cli *client.Client, c *Container) (bool, error) {
	summary, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return false, fmt.Errorf("Error while retrieving image list.\n%v", err)
	}

	for _, img := range summary {
		fmt.Println(c.ContainerConfig.Image)
		if found := contains(img.RepoTags, c.ContainerConfig.Image); found {
			return true, nil
		}
	}

	return false, nil
}

// pullImage pulls Docker image from specified repository
func pullImage(ctx context.Context, cli *client.Client, repo string) error {
	reader, err := cli.ImagePull(ctx, repo, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("Error while pulling image.\n%v", err)
	}
	io.Copy(os.Stdout, reader)

	return nil
}

// execCommand creates an exec configuration and runs an exec process
func execCommand(ctx context.Context, cli *client.Client, c *Container, cID ID) error {
	exec, err := cli.ContainerExecCreate(ctx, string(cID), *c.ExecConfig)
	if err != nil {
		return fmt.Errorf("Error while creating exec configuration.\n%v", err)
	}

	res, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return fmt.Errorf("Error while executing command.\n%v", err)
	}

	io.Copy(os.Stdout, res.Reader)

	return nil
}

// copyToContainer copies files into the container
func copyToContainer(ctx context.Context, cli *client.Client, c *Container, cID ID) error {
	for _, path := range c.Attachments {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error while opening file %s.\n%v", path, err)
		}
		tar, _ := archive.Tar(file.Name(), archive.Uncompressed)
		err = cli.CopyToContainer(ctx, string(cID), "/tmp/", tar, types.CopyToContainerOptions{
			AllowOverwriteDirWithFile: true,
		})
		if err != nil {
			return fmt.Errorf("Error while copying files to container.\n%v", err)
		}
		file.Close()
	}

	return nil
}

// disposeContainer stops and removes specified container
func disposeContainer(ctx context.Context, cli *client.Client, cID ID) {
	cli.ContainerStop(ctx, string(cID), nil)
	cli.ContainerRemove(ctx, string(cID), types.ContainerRemoveOptions{
		RemoveVolumes: true,
	})
}

// createContainer creates a container with specified configuration
func createContainer(ctx context.Context, cli *client.Client, c *Container) (ID, error) {
	cont, err := cli.ContainerCreate(ctx, c.ContainerConfig, c.HostConfig, nil, "")
	if err != nil {
		return "", fmt.Errorf("Error while creating the container.\n%v", err)
	}

	return ID(cont.ID), nil
}

// startContainer starts the specified container
func startContainer(ctx context.Context, cli *client.Client, cID ID) error {
	err := cli.ContainerStart(ctx, string(cID), types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("Error while starting container.\n%v", err)
	}

	fmt.Printf("Container %s started\n", cID)

	return nil
}

// RunContainer creates a new Container for the specified
// image (repo:tag) and sets host port and container port mappings (hp:cp).
// It downloads the image if it isn't found in Docker host.
func (c *Container) RunContainer() error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	if found, _ := imageExists(ctx, cli, c); found {
		goto createContainer
	}

	err = pullImage(ctx, cli, c.RemoteURL)
	if err != nil {
		return err
	}

createContainer:
	cID, err := createContainer(ctx, cli, c)
	if err != nil {
		return err
	}

	err = copyToContainer(ctx, cli, c, cID)
	if err != nil {
		return err
	}

	err = startContainer(ctx, cli, cID)
	if err != nil {
		return err
	}

	err = execCommand(ctx, cli, c, cID)
	if err != nil {
		return err
	}

	defer disposeContainer(ctx, cli, cID)

	return nil
}
