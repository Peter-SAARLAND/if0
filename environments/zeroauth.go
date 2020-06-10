package environments

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"math/rand"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandSeq() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateHashCmd(seq string) (string, error) {
	// htpasswd -nbB admin 'superpassword' | cut -d ":" -f 2
	if runtime.GOOS != "windows" {
		command := "htpasswd -nbB admin " + "'"+seq+"'" + "| cut -d \":\" -f 2"
		cmd := exec.Command("bash", "-c", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error: htpasswd - ", err, string(out))
			return "", err
		}
		hash := string(out)
		hash = strings.Replace(hash, "$", "$$", -1)
		return hash, nil
	}
	return "", nil
}

func generateHashDocker(seq string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	command := "htpasswd -nbB admin " + "'"+seq+"'" + "| cut -d \":\" -f 2"
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "httpd:2.4-alpine",
		Cmd:   []string{"bash", "-c", command},
		Tty:   true,
		AttachStdin: true,
		AttachStdout: true,
		AttachStderr: true,
	}, nil, nil, "htpwd")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, out)
	hash := buf.String()
	hash = strings.Replace(hash, "$", "$$", -1)

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return "", err
	}
	return hash, nil
}