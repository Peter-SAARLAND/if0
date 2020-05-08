package environments

import (
	"errors"
	"fmt"
	"if0/common"
	"if0/common/sync"
	"if0/config"
	"os"
	"path/filepath"
)

var (
	syncObj = sync.Sync{}
	clone = syncObj.Clone
	getAuth = sync.GetSyncAuth
	repoSync = config.GitRepoSync
)

func AddEnv(repoUrl string) error {
	// get authorization
	authObj := sync.Auth{}
	auth, err := getAuth(&authObj, repoUrl)
	if err != nil {
		fmt.Println("Authentication error - ", err)
		return err
	}

	_, err = clone(repoUrl, auth)
	if err != nil{
		return err
	}
	return nil
}

func SyncEnv(repoName string) error {
	repoPath := filepath.Join(common.EnvDir, repoName)
	// check if repo exists.
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("The repository %s could not be found at %s. " +
			"Add the repository before performing sync operation \n", repoName, common.EnvDir)
		return errors.New("repository not found")
	}

	err := repoSync(&syncObj, repoPath, false)
	if err != nil {
		fmt.Println("Error: Syncing external repo - ", err)
		return err
	}
	return nil
}