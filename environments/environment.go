package environments

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"if0/common"
	"if0/common/sync"
	"if0/config"
	"os"
	"path/filepath"
)

func AddEnv(repoUrl string) error {
	// get authorization
	authObj := sync.Auth{}
	auth, err := sync.GetSyncAuth(&authObj, repoUrl)
	if err != nil {
		log.Errorln("Authentication error: ", err)
		return err
	}

	syncObj := sync.Sync{}
	_, err = syncObj.Clone(repoUrl, auth)
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
		return errors.New("repo not found")
	}

	syncObj := sync.Sync{}
	err := config.GitRepoSync(&syncObj, repoPath, false)
	if err != nil {
		log.Errorln("Error while syncing external repo: ", err)
		return err
	}

	return nil
}