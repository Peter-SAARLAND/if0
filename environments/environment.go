package environments

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"if0/common"
	"if0/common/sync"
	"if0/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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
		if err.Error() == "remote repository is empty"{
			err = cloneEmptyRepo(repoUrl)
			if err != nil {
				return err
			}
		}
		return err
	}
	return nil
}

func cloneEmptyRepo(remoteStorage string) error {
	syncObj := sync.Sync{}
	dirName := strings.Split(filepath.Base(remoteStorage), ".")[0]
	dirPath := filepath.Join(common.EnvDir, dirName)
	r, err := syncObj.GitInit(dirPath)
	if err != nil {
		return err
	}
	// git remote add <repo>
	err = syncObj.AddRemote(remoteStorage, r)
	if err != nil {
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

func LoadEnv(envDir string) error {
	fmt.Println("Reading .env files from", envDir)
	envConfig := readAllEnvFiles(envDir)
	for k, v := range envConfig {
		config.SetEnvVariable(k, v.(string))
	}
	err := syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())
	if err != nil {
		fmt.Println("Error: Opening shell - ", err)
	}
	return nil
}

func readAllEnvFiles(dirPath string) map[string]interface{}{
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error: Reading environment directory %s - %s\n", dirPath, err)
		return nil
	}
	if len(files) < 1 {
		fmt.Println("Info: No .env files found")
		return nil
	}
	allConfig := make(map[string]interface{})
	for _, file := range files {
		fileName := filepath.Join(dirPath, file.Name())
		if filepath.Ext(fileName) == ".env" {
			viper.SetConfigFile(fileName)
			err := viper.ReadInConfig()
			if err != nil {
				fmt.Printf("Error: while reading %s file - %s\n", fileName, err)
				continue
			}
			currConfig := viper.AllSettings()
			for k, v := range currConfig {
				allConfig[k] = v
			}
		}
	}
	return allConfig
}