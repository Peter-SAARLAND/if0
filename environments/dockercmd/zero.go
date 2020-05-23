package dockercmd

import (
"fmt"
"github.com/docker/docker/api/types/container"
"github.com/docker/docker/api/types/mount"
)

// This function used to provision the platform
func MakeProvision(envName string) error {
	//binding mounts
	mounts := []mount.Mount{
		{Type: mount.TypeBind,
			Source: getMountSrcPath(envName),
			Target: mountTargetPath}}
	hostConfig := &container.HostConfig{Mounts: mounts}

	containerConfig := &container.Config{
		Image: zeroImage,
		Cmd:   []string{"make", "provision"},
		Tty:   true,
		Env:   []string{"IF0_ENVIRONMENT=" + envName},
	}
	containerName := "zero-" + envName
	err := dockerRun(containerConfig, hostConfig, containerName, zeroImage)
	if err != nil {
		fmt.Println("Error: MakeProvision - ", err)
		return err
	}
	return nil
}
