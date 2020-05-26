package environments

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"if0/common"
	"if0/common/sync"
	"if0/config"
	"if0/environments/dockercmd"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	syncObj  = sync.Sync{}
	clone    = syncObj.Clone
	getAuth  = sync.GetSyncAuth
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
	if err != nil {
		if err.Error() == "remote repository is empty" {
			err = cloneEmptyRepo(repoUrl)
			if err != nil {
				return err
			}
		}
		return err
	}

	//check if the necessary files are present in the environment
	//if not, add them with basic information
	envName := strings.Split(filepath.Base(repoUrl), ".")[0]
	envInit(envName)
	return nil
}

func SyncEnv(envDir string) error {
	// check if repo exists.
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		fmt.Printf("The repository could not be found at %s. "+
			"Please add the repository before performing sync operation \n", common.EnvDir)
		return errors.New("repository not found")
	}

	err := repoSync(&syncObj, envDir, false)
	if err != nil {
		fmt.Println("Error: Syncing external repo - ", err)
		return err
	}
	return nil
}

func PlanEnv(envDir string) error {
	envName := strings.Split(filepath.Base(envDir), ".")[0]
	err := dockercmd.MakePlan(envName)
	if err != nil {
		return err
	}
	return nil
}

func ProvisionEnv(envDir string) error {
	envName := strings.Split(filepath.Base(envDir), ".")[0]
	err := dockercmd.MakeProvision(envName)
	if err != nil {
		return err
	}
	return nil
}

func CreateZeroInfra(envDir string) error {
	envName := strings.Split(filepath.Base(envDir), ".")[0]
	err := dockercmd.MakeZero(envName)
	if err != nil {
		return err
	}
	return nil
}


func loadEnv(envDir string) error {
	fmt.Println("Reading .env files from", envDir)
	envConfig, err := readAllEnvFiles(envDir)
	if err != nil {
		return err
	}
	for k, v := range envConfig {
		config.SetEnvVariable(k, v.(string))
	}
	return nil
}

func readAllEnvFiles(dirPath string) (map[string]interface{}, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error: Reading environment directory %s - %s\n", dirPath, err)
		return nil, err
	}
	if len(files) < 1 {
		fmt.Println("Info: No .env files found")
		return nil, errors.New("no .env files found")
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
	return allConfig, nil
}