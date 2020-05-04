package environments

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"if0/common"
	"os"
	"path/filepath"
)

func AddEnv(repoUrl string) error {
	// get authorization
	authObj := common.Auth{}
	auth, err := common.GetSyncAuth(&authObj, repoUrl)
	if err != nil {
		log.Errorln("Authentication error: ", err)
		return err
	}

	syncObj := common.Sync{}
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
			"Add the repository before performing sync operation", repoName, common.EnvDir)
		return errors.New("repo not found")
	}
	
	return nil
}