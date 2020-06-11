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
	getAuth  = sync.GetSyncAuth
	repoSync = config.GitRepoSync
)

func AddEnv(addEnvArgs []string) error {
	var repoName, repoUrl string
	repoName = addEnvArgs[0]
	if len(addEnvArgs) > 1 {
		repoUrl = addEnvArgs[1]
	}
	config.ReadConfigFile(common.If0Default)
	gitlabToken := config.GetEnvVariable("GL_TOKEN")
	if gitlabToken == "" || repoUrl != "" {
		// adding environment locally (to sync with later)
		// or syncing a local environment that has already been added
		err := createLocalEnv(repoName, repoUrl)
		if err != nil {
			return err
		}
	} else {
		// adding environment using GitLab token
		err := createGLProject(repoName, gitlabToken)
		if err != nil {
			fmt.Println("Error: Adding Private Project -", err)
			return err
		}
	}
	return nil
}

func SyncEnv(envDir string) error {
	// check if repo exists.
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		fmt.Printf("The repository could not be found at %s. "+
			"Please add the repository before performing sync operation \n", common.EnvDir)
		return errors.New("repository not found")
	}

	err := repoSync(&syncObj, "", envDir)
	if err != nil {
		fmt.Println("Error: Syncing external repo - ", err)
		return err
	}
	return nil
}

func Dash1Plan(envDir string) error {
	envName := strings.Split(filepath.Base(envDir), ".")[0]
	err := dockercmd.MakePlan(envName)
	if err != nil {
		return err
	}
	return nil
}

func ZeroPlatform(envDir string) error {
	envName := strings.Split(filepath.Base(envDir), ".")[0]
	err := dockercmd.MakePlatform(envName)
	if err != nil {
		return err
	}
	return nil
}

func Dash1Infrastructure(envDir string) error {
	envName := strings.Split(filepath.Base(envDir), ".")[0]
	err := dockercmd.MakeInfrastructure(envName)
	if err != nil {
		return err
	}
	return nil
}

func Dash1Destroy(envDir string) error {
	envName := strings.Split(filepath.Base(envDir), ".")[0]
	err := dockercmd.MakeDestroy(envName)
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
