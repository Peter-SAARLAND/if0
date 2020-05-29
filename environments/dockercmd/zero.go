package dockercmd

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
)

// This function used to provision the platform
func MakeProvision(envName string) error {
	//binding mounts
	mounts := addMounts(envName)
	if mounts == nil {
		errString := fmt.Sprintf("environment %s doesn't exist. "+
			"Do `if0 environment add %s` to add it", envName, envName)
		return errors.New(errString)
	}
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
