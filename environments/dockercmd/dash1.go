package dockercmd

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"strings"
)

// This function is used to start a dash1 container, and run `make plan` inside the container.
// In dash1, make plan initializes the necessary Terraform provider modules for
// the Environment 'envName' and then creates a plan in ~/.if0/.environments/$NAME/dash1.plan`
func MakePlan(envName string) error {
	command := []string{"make", "plan"}
	return dash1make(envName, command)
}

func MakeInfrastructure(envName string) error {
	command := []string{"make", "infrastructure"}
	return dash1make(envName, command)
}

func MakeDestroy(envName string) error {
	command := []string{"make", "destroy"}
	return dash1make(envName, command)
}

func dash1make(envName string, command []string) error {
	//binding mounts
	mounts := addMounts(envName)
	if mounts == nil {
		errString := fmt.Sprintf("environment %s doesn't exist. "+
			"Do `if0 environment add %s` to add it", envName, envName)
		return errors.New(errString)
	}
	hostConfig := &container.HostConfig{Mounts: mounts}

	containerConfig := &container.Config{
		Image: dash1Image,
		Cmd:   command,
		Tty:   true,
		Env:   []string{"IF0_ENVIRONMENT=" + envName},
	}
	envSplit := strings.Split(envName, "/")
	env := envSplit[len(envSplit)-1]
	containerName := "dash1-" + env
	err := dockerRun(containerConfig, hostConfig, containerName, dash1Image)
	if err != nil {
		return err
	}
	return nil
}
