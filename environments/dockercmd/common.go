package dockercmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"if0/common"
	"io"
	"os"
	"path/filepath"
)

const (
	dash1Image      = "registry.gitlab.com/peter.saarland/dash1"
	zeroImage       = "registry.gitlab.com/peter.saarland/zero"
	mountTargetPath = "/root/.if0/.environments/zero"
	gitConfigTargetPath = "/root/.gitconfig"
)

var (
	gitConfigSrc = filepath.Join(common.RootPath, ".gitconfig")
)

func getMountSrcPath(envName string) string {
	return filepath.Join(common.EnvDir, envName)
}

func dockerRun(containerConfig *container.Config, hostConfig *container.HostConfig,
	containerName string, image string) error {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error: ContainerClient - ", err)
		return err
	}

	status, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		fmt.Println("Error: ImagePull - ", err)
		return err
	}
	_, _ = io.Copy(os.Stdout, status)

	resp, err := dockerClient.ContainerCreate(ctx, containerConfig,
		hostConfig, nil, containerName)
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
