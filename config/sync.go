package config

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func RepoSync() error {
	fmt.Println("all: ", viper.AllSettings())
	remoteStorage := GetEnvVariable("REMOTE_STORAGE")
	fmt.Println("remoteStorage: ", remoteStorage)
	if remoteStorage == "" {
		return errors.New("REMOTE_STORAGE is not set.")
	}
	err := gitSync(remoteStorage)
	if err != nil {
		log.Errorln("Error while syncing external repo: ", err)
		return err
	}
	return nil
}

func gitSync(remoteStorage string) error {
	// check if the repo is already present
	repoDir := filepath.Join(if0Dir, "if0")
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		log.Debugln("Directory does not exist, creating dir for if0 repository")
		_ = os.Mkdir(repoDir, os.ModeDir)
		if strings.Contains(remoteStorage, "http") {
			err := gitHttpsClone(remoteStorage, repoDir)
			if err != nil {
				return err
			}
		} else if strings.Contains(remoteStorage, "git@") {
			err := gitSSHClone(remoteStorage, repoDir)
			if err != nil {
				return err
			}
		}
	} else {
		// git pull logic
		var auth transport.AuthMethod
		if strings.Contains(remoteStorage, "http") {
			auth, err = getHttpAuthCredentials()
			if err != nil {
				log.Errorln("Error while fetching credentials: ", err)
				return err
			}
		} else if strings.Contains(remoteStorage, "git@") {
			auth, err = getSSHAuth()
			if err != nil {
				log.Errorln("Error while fetching credentials: ", err)
				return err
			}
		}
		gitPull(auth)
	}
	return nil
}

func gitPull(auth transport.AuthMethod) {
	r, err := git.PlainOpen(filepath.Join(if0Dir, "if0"))
	if err != nil {
		log.Println(err)
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		log.Println(err)
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	log.Println("git pull origin")
	err = w.Pull(&git.PullOptions{RemoteName: "origin", Auth: auth})
	if err != nil {
		log.Println(err)
	}
}

func gitSSHClone(remoteStorage string, repoDir string) error {
	auth, err := getSSHAuth()
	if err != nil {
		return err
	}
	_, err = gitClone(remoteStorage, auth, repoDir)
	if err != nil {
		return err
	}
	return nil
}

func getSSHAuth() (*gitssh.PublicKeys, error) {
	sshKeyPath := GetEnvVariable("SSH_KEY_PATH")
	sshKey, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		fmt.Println("ssh err: ", err)
		return nil, err
	}
	passphrase := getPassphrase()
	signer, err := ssh.ParsePrivateKeyWithPassphrase(sshKey, passphrase)
	if err != nil {
		fmt.Println("signer err: ", err)
		return nil, err
	}
	auth := &gitssh.PublicKeys{User: "git", Signer: signer}
	return auth, nil
}

func gitHttpsClone(remoteStorage string, repoDir string) error {
	auth, err := getHttpAuthCredentials()
	if err != nil {
		log.Errorln("Error while fetching credentials: ", err)
		return err
	}
	log.Printf("Cloning from %s at %s\n", remoteStorage, repoDir)
	_, err = gitClone(remoteStorage, auth, repoDir)
	if err != nil {
		return err
	}
	return nil
}

func gitClone(remoteStorage string, auth transport.AuthMethod, dir string) (*git.Repository, error) {
	cloneOptions := &git.CloneOptions{URL: remoteStorage,
		Auth: auth, Progress: os.Stdout}
	r, err := git.PlainClone(dir, false, cloneOptions)
	if err != nil {
		log.Errorf("Error while cloning the repo from %s - %s\n", remoteStorage, err)
		return nil, err
	}
	return r, nil
}

func getHttpAuthCredentials() (transport.AuthMethod, error) {
	fmt.Println("Enter Username: ")
	userName, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Failed to read username: %v", err)
	}
	fmt.Println("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Failed to read password: %v", err)
	}
	auth := &http.BasicAuth{Username: string(userName), Password: string(bytePassword)}
	return auth, err
}

func getPassphrase() []byte {
	fmt.Println("Enter Passphrase. If you do not have a passphrase, press enter.")
	passphrase, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Failed to read username: %v", err)
	}
	return []byte(passphrase)
}
