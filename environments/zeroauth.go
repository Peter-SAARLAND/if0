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
	command := "htpasswd -nbB admin " + "'"+seq+"'" + "| cut -d \":\" -f 2"
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err", err, string(out))
		return "", err
	}
	hash := string(out)
	hash = strings.Replace(hash, "$", "$$", -1)
	return hash, nil
}

func generateHashDocker(seq string) string {
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
		fmt.Println("Error: ContainerCreate - ", err)
		return ""
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("Error: ContainerStart - ", err)
		return ""
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Error: ContainerWait - ", err)
			return ""
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		fmt.Println("Error: ContainerLogs - ", err)
		return ""
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, out)
	hash := buf.String()
	hash = strings.Replace(hash, "$", "$$", -1)

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		fmt.Println("Error: ContainerRemove - ", err)
		return ""
	}
	return hash
}