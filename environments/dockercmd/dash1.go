package dockercmd

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

// This function is used to start a dash1 container, and run `make plan` inside the container.
// In dash1, make plan initializes the necessary Terraform provider modules for
// the Environment 'envName' and then creates a plan in ~/.if0/.environments/$NAME/dash1.plan`
func MakePlan(envName string) error {
	//binding mounts
	mounts := []mount.Mount{
		{Type: mount.TypeBind,
			Source: getMountSrcPath(envName),
			Target: mountTargetPath},
		{Type: mount.TypeBind,
			Source: gitConfigSrc,
			Target: gitConfigTargetPath}}
	hostConfig := &container.HostConfig{Mounts: mounts}

	containerConfig := &container.Config{
		Image: dash1Image,
		Cmd:   []string{"make", "plan"},
		Tty:   true,
		Env:   []string{"IF0_ENVIRONMENT=" + envName},
	}
	containerName := "dash1-" + envName
	err := dockerRun(containerConfig, hostConfig, containerName, dash1Image)
	if err != nil {
		fmt.Println("Error: MakePlan - ", err)
		return err
	}
	return nil
}

// This function used to provision the platform
func MakeZero(envName string) error {
	//binding mounts
	mounts := []mount.Mount{
		{Type: mount.TypeBind,
			Source: getMountSrcPath(envName),
			Target: mountTargetPath}}
	hostConfig := &container.HostConfig{Mounts: mounts}

	containerConfig := &container.Config{
		Image: dash1Image,
		Cmd:   []string{"make", "zero"},
		Tty:   true,
		Env:   []string{"IF0_ENVIRONMENT=" + envName},
	}
	containerName := "dash1-" + envName
	err := dockerRun(containerConfig, hostConfig, containerName, dash1Image)
	if err != nil {
		fmt.Println("Error: MakeZero - ", err)
		return err
	}
	return nil
}
