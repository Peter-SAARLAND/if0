package config

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/pkg/errors"
	"if0/common"
	"if0/common/sync"
	"os"
	"path/filepath"
	"strings"
)

var (
	GitRepoSync          = GitSync
	repoUrl              = getRepoUrl
	checkForLocalChanges = localChanges
)

// RepoSync is used to synchronize the if0 configuration files with a remote git repository
func RepoSync() error {
	remoteStorage := GetEnvVariable("REMOTE_STORAGE")
	fmt.Println("Syncing with remote storage: ", remoteStorage)
	if remoteStorage == "" {
		return errors.New("REMOTE_STORAGE is not set.")
	}
	syncObj := sync.Sync{}
	err := GitRepoSync(&syncObj, remoteStorage, common.If0Dir)
	if err != nil {
		fmt.Println("Error:Syncing external repo - ", err)
		return err
	}
	return nil
}

func GitSync(syncObj sync.SyncOps, repo string, dir string) error {
	// get repository (git init, remote add; or open an existing repository)
	r, err := GetRepository(syncObj, repo, dir)
	if err != nil {
		return err
	}
	repo = repoUrl(r)

	// get authorization
	// if the git sync is via HTTPS, then fetch username-password credentials
	// if the git sync is via SSH, then parse .ppk file
	authObj := sync.Auth{}
	auth, err := sync.GetSyncAuth(&authObj, repo)
	if err != nil {
		fmt.Println("Authentication Error - ", err)
		return err
	}

	auto, manual, err := checkForLocalChanges(syncObj, r)
	if err != nil {
		return err
	}

	if manual {
		fmt.Println("Exiting.")
		return errors.New("add/commit the local changes before sync")
	}

	pullOptions := &git.PullOptions{Auth: auth, RemoteName: "origin", Force: false}
	_, err = syncObj.Pull(repo, r, pullOptions)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate || err.Error() == "remote repository is empty" {
			fmt.Println("Pull status: ", err)
		} else {
			fmt.Println("Error: Pull status - ", err)
			return err
		}
	}

	if auto {
		err = syncChanges(syncObj, r, auth)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRepository(syncObj sync.SyncOps, repo string, dir string) (*git.Repository, error) {
	r := &git.Repository{}
	// check if the repo is already present
	// if not, do a `git init`, and `git remote add origin repo`
	if _, err := os.Stat(filepath.Join(dir, git.GitDirName)); os.IsNotExist(err) && repo != "" {
		// git init
		r, err = syncObj.GitInit(dir)
		if err != nil {
			return nil, err
		}

		// git remote add <repo>
		err = syncObj.AddRemote(repo, r)
		if err != nil {
			return nil, err
		}
	} else {
		// open the existing repo at ~/.if0
		r, err = syncObj.Open(dir)
		if err != nil {
			fmt.Println("Error: Opening repository - ", err)
			return nil, err
		}
	}
	return r, nil
}

func syncChanges(syncObj sync.SyncOps, r *git.Repository, auth transport.AuthMethod) error {
	fmt.Println("Pushing the local changes")
	w, err := syncObj.GetWorktree(r)
	if err != nil {
		fmt.Println("Worktree Error: ", err)
	}
	// git commit
	err = syncObj.Commit(w)
	if err != nil {
		fmt.Println("Error: Committing changes - ", err)
	}
	// git push
	err = syncObj.Push(auth, r)
	if err != nil {
		fmt.Println("Error: Pushing changes - ", err)
		return err
	}
	return nil
}

func getRepoUrl(r *git.Repository) string {
	remotes, err := r.Remote("origin")
	if err != nil {
		fmt.Println("Error: Remotes - ", err)
	}
	return remotes.Config().URLs[0]
}

func localChanges(syncObj sync.SyncOps, r *git.Repository) (bool, bool, error) {
	var auto, manual bool
	// worktree
	w, err := syncObj.GetWorktree(r)
	if err != nil {
		fmt.Println("Worktree Error - ", err)
		return false, false, err
	}

	status, err := getStatus(syncObj, w)
	if err != nil {
		return false, false, err
	}

	if len(status) > 0 {
		// prompt the user if they want to add/commit/push changes
		fmt.Println("Following changes were found. " +
			"If the repository is not up-to-date, pulling in changes would delete the unstaged changes. \n" +
			"Other changes would be overwritten by the remote changes.")
		fmt.Println(status)
		fmt.Println("Proceed? [Y/n]")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// if the user chooses 'y', add/commit/push changes
		switch strings.TrimSpace(strings.ToLower(text)) {
		case "y", "":
			auto = true
			manual = false
			fmt.Println("Adding local changes")
			for file, _ := range status {
				err := syncObj.AddFile(w, file)
				if err != nil {
					fmt.Printf("Error: Adding file %s: %s \n", file, err)
					return false, false, err
				}
			}
		case "n":
			auto = false
			manual = true
		}
	}
	return auto, manual, nil
}

func getStatus(syncObj sync.SyncOps, w *git.Worktree) (git.Status, error) {
	// git status
	status, err := syncObj.Status(w)
	if err != nil {
		fmt.Println("Error: Status - ", err)
		return nil, err
	}
	return status, nil
}
