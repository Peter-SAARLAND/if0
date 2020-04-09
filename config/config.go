package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	if0Default = "if0.env"
)

var (
	rootPath, _  = os.Getwd() // TODO: change this to root dir ~
	envDir       = filepath.Join(rootPath, ".environments")
	snapshotsDir = filepath.Join(rootPath, ".snapshots")
)

func SetEnvVariable(key, value string) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	viper.Set(key, value)
	if value == GetEnvVariable(key) {
		log.Printf("key %s update with %s successful \n", key, value)
	}
}

func GetEnvVariable(key string) string {
	val := viper.Get(key)
	return cast.ToString(val)
}

// PrintCurrentRunningConfig reads the current running if0/env configuration file and prints it
func PrintCurrentRunningConfig() {
	readConfigFile(filepath.Join(rootPath, if0Default))
	allConfig := viper.AllSettings()
	for key, val := range allConfig {
		fmt.Println(key, ": ", val)
	}
}

// AddConfigFile replaces the current config file with the provided config file.
// it first checks if the file is already present,
// if present, it creates a backup in the ~if0/.snapshots directory
// and then proceeds to replace the current running config file
func AddConfigFile(srcConfigFile string, zero, merge bool) {
	runningConfigFile := getRunningConfigFile(srcConfigFile, zero)
	// taking a backup of the running configuration if already present
	present := isFilePresent(runningConfigFile, zero)
	if present {
		err := backupToSnapshots(runningConfigFile)
		if err != nil {
			log.Fatalln("Failed to add/update the config file: ", err)
		}
	}
	if merge {
		mergeConfigFiles(srcConfigFile, runningConfigFile)
	} else {
		createConfigFile(srcConfigFile, runningConfigFile)
	}
}