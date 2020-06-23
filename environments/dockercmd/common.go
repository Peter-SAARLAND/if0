package dockercmd

import (
	"context"
	"errors"
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
	dash1Image          = "registry.gitlab.com/peter.saarland/dash1"
	zeroImage           = "registry.gitlab.com/peter.saarland/zero"
	mountTargetPath     = "/root/.if0/.environments/zero"
	gitConfigTargetPath = "/root/.gitconfig"
)

func addMounts(envName string) []mount.Mount {
	var mounts []mount.Mount
	mountPath, err := getMountSrcPath(envName)
	if err != nil {
		return nil
	}
	mounts = append(mounts, mount.Mount{Type: mount.TypeBind,
		Source: mountPath,
		Target: mountTargetPath})
	// append gitconfig mount, if present.
	mounts = getGitConfigMount(mounts)
	return mounts
}

func getMountSrcPath(envName string) (string, error) {
	envDir := filepath.Join(common.EnvDir, envName)
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return "", errors.New("missing mount")
	}
	return filepath.Join(common.EnvDir, envName), nil
}

func getGitConfigMount(mounts []mount.Mount) []mount.Mount {
	gitConfigPath := getGitConfigPath()
	if gitConfigPath != "" {
		mounts = append(mounts, mount.Mount{Type: mount.TypeBind,
			Source: gitConfigPath, Target: gitConfigTargetPath})
	}
	return mounts
}

func getGitConfigPath() string {
	gitConfigSrc := filepath.Join(common.RootPath, ".gitconfig")
	if _, err := os.Stat(gitConfigSrc); os.IsNotExist(err) {
		return ""
	}
	return gitConfigSrc
}

func dockerRun(containerConfig *container.Config, hostConfig *container.HostConfig,
	containerName string, image string) error {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error: ContainerClient -", err)
		return err
	}

	containerConfig.Env = append(containerConfig.Env, "VERBOSITY=1")
	resp, err := dockerClient.ContainerCreate(ctx, containerConfig,
		hostConfig, nil, containerName)
	if err != nil {
		fmt.Println("Error: ContainerCreate -", err)
		return err
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("Error: ContainerStart -", err)
		_ = removeContainer(dockerClient, resp.ID)
		return err
	}

	printContainerLogs(err, dockerClient, ctx, resp)

	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Error: ContainerWait -", err)
			_ = removeContainer(dockerClient, resp.ID)
			return err
		}
	case <-statusCh:
	}

	err = removeContainer(dockerClient, resp.ID)
	if err != nil {
		return err
	}
	return nil
}

func printContainerLogs(err error, dockerClient *client.Client, ctx context.Context, resp container.ContainerCreateCreatedBody) {
	out, err := dockerClient.ContainerLogs(ctx, resp.ID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
	if err != nil {
		fmt.Println("Error: ContainerLogs -", err)
	}
	defer out.Close()
	_, _ = io.Copy(os.Stdout, out)
}

func removeContainer(dockerClient *client.Client, respId string) error {
	ctx := context.Background()
	err := dockerClient.ContainerRemove(ctx, respId, types.ContainerRemoveOptions{})
	if err != nil {
		fmt.Println("Error: ContainerRemove -", err)
		return err
	}
	return nil
}
