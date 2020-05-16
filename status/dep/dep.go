package dep

import (
	"fmt"
	"os/exec"
	"runtime"
)

func CheckDependencies() {
	fmt.Println("Checking docker dependency...")
	dependency := "docker"
	version, err := checkDependencyVersion(dependency)
	if err != nil {
		fmt.Println("Error checking for docker version - ", err)
	} else {
		fmt.Println(version)
	}

	fmt.Println("Checking vagrant dependency...")
	dependency = "vagrant"
	version, err = checkDependencyVersion(dependency)
	if err != nil {
		fmt.Println("Error checking for vagrant version", err)
	} else {
		fmt.Println(version)
	}
}

func checkDependencyVersion(dependency string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", dependency, "--version")
	} else {
		cmd = exec.Command(dependency, "--version")
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}