package dockercmd

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"strings"
)

// This function used to provision the platform
func MakePlatform(envName string) error {
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
		Cmd:   []string{"make", "platform"},
		Tty:   true,
		Env:   []string{"IF0_ENVIRONMENT=" + envName},
	}
	envSplit := strings.Split(envName, "/")
	env := envSplit[len(envSplit)-1]
	containerName := "zero-" + env
	err := dockerRun(containerConfig, hostConfig, containerName, zeroImage)
	if err != nil {
		fmt.Println("Error: MakePlatform - ", err)
		return err
	}
	return nil
}
