package config

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

var (
	repoSync = gitSync
)

// RepoSync is used to synchronize the if0 configuration files with a remote git repository
func RepoSync() error {
	remoteStorage := GetEnvVariable("REMOTE_STORAGE")
	log.Println("Syncing with remote storage: ", remoteStorage)
	if remoteStorage == "" {
		return errors.New("REMOTE_STORAGE is not set.")
	}
	syncObj := sync{}
	err := repoSync(&syncObj, remoteStorage)
	if err != nil {
		log.Errorln("Error while syncing external repo: ", err)
		return err
	}
	return nil
}

func gitSync(syncObj syncOps, remoteStorage string) error {
	// get authorization
	// if the git sync is via HTTPS, then fetch username-password credentials
	// if the git sync is via SSH, then parse .ppk file
	authObj := auth{}
	auth, err := getSyncAuth(&authObj, remoteStorage)
	if err != nil {
		log.Errorln("Authentication error: ", err)
		return err
	}

	// check if the repo is already present
	// if not, do a `git init`, and `git remote add origin remoteStorage`
	r := &git.Repository{}
	if _, err := os.Stat(filepath.Join(if0Dir, git.GitDirName)); os.IsNotExist(err) {
		// git init
		r, err = syncObj.gitInit(if0Dir)
		if err != nil {
			return err
		}

		// git remote add <repo>
		err = syncObj.addRemote(remoteStorage, r)
		if err != nil {
			return err
		}
	} else {
		// open the existing repo at ~/.if0
		r, err = syncObj.open(if0Dir)
		if err != nil {
			log.Errorln("Error while opening repository: ", err)
			return err
		}
	}

	// git pull
	pullOptions := &git.PullOptions{Auth: auth, RemoteName: "origin"}
	worktree, err := syncObj.pull(remoteStorage, r, pullOptions)
	if err != nil {
		log.Errorln("Pull status: ", err)
	}

	fmt.Println("worktree ", worktree)
	// git status
	status, err := syncObj.status(worktree)
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
		err := syncObj.addFile(worktree, file)
		if err != nil {
			log.Errorf("Error adding file %s: %s \n", file, err)
		}
	}

	// git commit
	err = syncObj.commit(worktree)
	if err != nil {
		log.Errorln("Error while committing changes: ", err)
		return err
	}

	// git push
	err = syncObj.push(auth, r)
	if err != nil {
		log.Errorln("Error while pushing changes: ", err)
		return err
	}
	return nil
}