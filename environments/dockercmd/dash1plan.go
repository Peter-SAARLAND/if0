package dockercmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"if0/common"
	"io"
	"os"
	"path/filepath"
)

const (
	dash1Image = "registry.gitlab.com/peter.saarland/dash1"
)

// This function is used to start a dash1 container, and run `make plan` inside the container.
// In dash1, make plan initializes the necessary Terraform provider modules for
// the Environment 'envName' and then creates a plan in ~/.if0/.environments/$NAME/dash1.plan`
func MakePlan(envName string) error {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error: ContainerClient - ", err)
		return err
	}

	status, err := dockerClient.ImagePull(ctx, dash1Image, types.ImagePullOptions{})
	if err != nil {
		fmt.Println("Error: ImagePull - ", err)
		return err
	}
	_, _ = io.Copy(os.Stdout, status)

	//binding mounts
	mounts := []mount.Mount{
		{Type: mount.TypeBind,
			Source: filepath.Join(common.RootPath, ".if0", ".environments", "if0-config"),
			Target: "/root/.if0/.environments/zero"},
		{Type: mount.TypeBind,
			Source: filepath.Join(common.RootPath, ".gitconfig"),
			Target: "/root/.gitconfig"}}
	hostConfig := &container.HostConfig{Mounts: mounts}

	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: dash1Image,
		Cmd:   []string{"make", "plan"},
		Tty:   true,
	}, hostConfig, nil, "dash1-"+envName)
	if err != nil {
		fmt.Println("Error: ContainerCreate - ", err)
		return err
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("Error: ContainerStart - ", err)
		return err
	}

	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Error: ContainerWait - ", err)
			return err
		}
	case <-statusCh:
	}

	out, err := dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		fmt.Println("Error: ContainerLogs - ", err)
		return err
	}
	_, _ = io.Copy(os.Stdout, out)

	err = dockerClient.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		fmt.Println("Error: ContainerRemove - ", err)
		return err
	}
	return nil
}
