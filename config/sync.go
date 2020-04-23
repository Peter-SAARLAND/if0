package config

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// RepoSync is used to synchronize the if0 configuration files with a remote git repository
func RepoSync() error {
	remoteStorage := GetEnvVariable("REMOTE_STORAGE")
	log.Println("Syncing with remote storage: ", remoteStorage)
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
	// get authorization
	// if the git sync is via HTTPS, then fetch username-password credentials
	// if the git sync is via SSH, then parse .ppk file
	auth, err := getAuth(remoteStorage)
	if err != nil {
		log.Errorln("Authentication error: ", err)
		return err
	}

	// check if the repo is already present
	// if not, do a `git init`, and `git remote add origin remoteStorage`
	r := &git.Repository{}
	if _, err := os.Stat(filepath.Join(if0Dir, git.GitDirName)); os.IsNotExist(err) {
		// git init
		r, err = gitInit(if0Dir, r)
		if err != nil {
			return err
		}

		// git remote add <repo>
		err = addRemote(remoteStorage, r)
		if err != nil {
			return err
		}
	} else {
		log.Println("Git repository already present.")
		// open the existing repo at ~/.if0
		r, err = git.PlainOpen(if0Dir)
		if err != nil {
			log.Errorln("Error while opening repository: ", err)
			return err
		}
	}

	// git pull
	log.Println("Pulling in changes from ", remoteStorage)
	w, err := r.Worktree()
	if err != nil {
		log.Errorln("worktree error: ", err)
		return err
	}
	pullOptions := &git.PullOptions{Auth: auth, RemoteName: "origin"}
	err = w.Pull(pullOptions)
	if err != nil {
		log.Errorln("Pull status: ", err)
	}

	// git status
	status, err := w.Status()
	if err != nil {
		fmt.Println("status err: ", err)
	}
	if len(status) == 0 {
		log.Println("No local changes were found. Exiting")
		return nil
	}

	// prompt the user if they want to add/commit/push changes
	fmt.Println("Following changes were found. Would you like to commit them?")
	fmt.Println(status)
	fmt.Println("Enter y or n")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// exit if the user does not want to add/commit/push changes
	if strings.EqualFold(strings.TrimSpace(text), "n") {
		log.Println("Exiting")
		return nil
	}

	// git add
	log.Println("Adding local changes")
	for file, _ := range status {
		_, err := w.Add(file)
		if err != nil {
			fmt.Printf("error adding file %s: %s \n", file, err)
		}
	}

	// git commit

	commitOptions := &git.CommitOptions{
		All: false,
		Author: &object.Signature{
			When:  time.Now(),
		},
	}
	commitMsg := "feat: updating config files - " + time.Now().Format("02012006_150405")
	_, err = w.Commit(commitMsg, commitOptions)
	if err != nil {
		log.Errorln("Error while committing changes: ", err)
		return err
	}

	// git push
	log.Println("Pushing local changes")
	pushOptions := &git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	}
	err = r.Push(pushOptions)
	if err != nil {
		log.Errorln("Error while pushing changes: ", err)
		return err
	}
	return nil
}

func getAuth(remoteStorage string) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	var err error
	if strings.Contains(remoteStorage, "http") {
		auth, err = getHttpAuth()
		if err != nil {
			log.Errorln("Error while fetching credentials: ", err)
			return nil, err
		}
	} else if strings.Contains(remoteStorage, "git@") {
		auth, err = getSSHAuth()
		if err != nil {
			log.Errorln("Error while fetching credentials: ", err)
			return nil, err
		}
	}
	return auth, nil
}

func getHttpAuth() (transport.AuthMethod, error) {
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

func getPassphrase() []byte {
	fmt.Println("Enter Passphrase. If you do not have a passphrase, press enter.")
	passphrase, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Failed to read username: %v", err)
	}
	return []byte(passphrase)
}

func gitInit(localRepoPath string, r *git.Repository) (*git.Repository, error) {
	log.Println("Creating a git repository at ", localRepoPath)
	r, err := git.PlainInit(localRepoPath, false)
	if err != nil {
		log.Errorln("Error while creating a git repository: ", err)
		return nil, err
	}
	return r, nil
}

func addRemote(remoteStorage string, r *git.Repository) error {
	log.Println("Adding remote 'origin' for the repository at ", remoteStorage)
	remoteConfig := &config.RemoteConfig{Name: "origin", URLs: []string{remoteStorage}}
	_, err := r.CreateRemote(remoteConfig)
	if err != nil {
		log.Errorln("Error while adding remote: ", err)
		return err
	}
	return nil
}
