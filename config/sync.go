package config

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"if0/common"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	syncObj := common.Sync{}
	err := repoSync(&syncObj, remoteStorage)
	if err != nil {
		log.Errorln("Error while syncing external repo: ", err)
		return err
	}
	return nil
}

func gitSync(syncObj common.SyncOps, remoteStorage string) error {
	// get authorization
	// if the git sync is via HTTPS, then fetch username-password credentials
	// if the git sync is via SSH, then parse .ppk file
	authObj := common.Auth{}
	auth, err := common.GetSyncAuth(&authObj, remoteStorage)
	if err != nil {
		log.Errorln("Authentication error: ", err)
		return err
	}

	// check if the repo is already present
	// if not, do a `git init`, and `git remote add origin remoteStorage`
	r := &git.Repository{}
	if _, err := os.Stat(filepath.Join(common.If0Dir, git.GitDirName)); os.IsNotExist(err) {
		// git init
		r, err = syncObj.GitInit(common.If0Dir)
		if err != nil {
			return err
		}

		// git remote add <repo>
		err = syncObj.AddRemote(remoteStorage, r)
		if err != nil {
			return err
		}
	} else {
		// open the existing repo at ~/.if0
		r, err = syncObj.Open(common.If0Dir)
		if err != nil {
			log.Errorln("Error while opening repository: ", err)
			return err
		}
	}

	// git pull
	//err = gitPull(if0Dir)
	//if err != nil {
	//	log.Errorln("Git pull error -", err)
	//	return err
	//}

	auto, manual, err := checkForLocalChanges(syncObj, r)
	if err != nil {
		return err
	}

	if manual {
		fmt.Println("Exiting.")
		return errors.New("add/commit the local changes before sync")
	}

	pullOptions := &git.PullOptions{Auth: auth, RemoteName: "origin", Force: false}
	_, err = syncObj.Pull(remoteStorage, r, pullOptions)
	if err != nil {
		log.Errorln("Pull status: ", err)
	}

	if auto {
		fmt.Println("Pushing the local changes")
		w, err := r.Worktree()
		if err != nil {
			log.Errorln("Worktree error: ", err)
		}
		// git commit
		err = syncObj.Commit(w)
		if err != nil {
			log.Errorln("Error while committing changes: ", err)
		}
		// git push
		err = syncObj.Push(auth, r)
		if err != nil {
			log.Errorln("Error while pushing changes: ", err)
			return err
		}
	}
	return nil
}

func checkForLocalChanges(syncObj common.SyncOps, r *git.Repository) (bool, bool, error) {
	var auto, manual bool
	// worktree
	w, err := r.Worktree()
	if err != nil {
		log.Errorln("Worktree error: ", err)
		return false, false, err
	}

	status, err := checkStatus(syncObj, w)
	if err != nil {
		return false, false, err
	}

	if len(status) > 0 {
		// prompt the user if they want to add/commit/push changes
		fmt.Println("Following changes were found. " +
			"Pulling in changes from the remote repository would " +
			"delete the unstaged changes before git init. \n" +
			"Other changes would be overwritten by the remote changes.")
		fmt.Println(status)
		fmt.Println("Enter 'y' to proceed. \n" +
			"Enter 'n' to exit `sync` and add the changes manually. ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// if the user chooses 'y', add/commit/push changes
		switch strings.TrimSpace(text) {
		case "y":
			auto = true
			manual = false
			log.Println("Adding local changes")
			for file, _ := range status {
				err := syncObj.AddFile(w, file)
				if err != nil {
					log.Errorf("Error adding file %s: %s \n", file, err)
					return false, false, err
				}
			}
		case "m":
			auto = false
			manual = true
		case "n":
			auto = false
			manual = false
		}
	}
	return auto, manual, nil
}

func checkStatus(syncObj common.SyncOps, w *git.Worktree) (git.Status, error) {
	// git status
	status, err := syncObj.Status(w)
	if err != nil {
		fmt.Println("status err: ", err)
		return nil, err
	}
	return status, nil
}

// WORKAROUND for git pull as git.Pull deletes unstaged changes.
func gitPull(if0Dir string) error {
	err := os.Chdir(if0Dir)
	if err != nil {
		fmt.Println("err chdir - ", err)
	}
	log.Println("Pulling changes")
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "git", "pull", "origin", "master")
	} else {
		cmd = exec.Command("git", "pull", "origin", "master")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorln("Error while doing git pull - ", string(out))
		return errors.New(string(out))
	}
	return nil
}