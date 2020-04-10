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
	rootPath, _  = os.UserHomeDir()
	if0Dir       = filepath.Join(rootPath, ".if0")
	envDir       = filepath.Join(if0Dir, ".environments")
	snapshotsDir = filepath.Join(if0Dir, ".snapshots")
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
	if0File := filepath.Join(if0Dir, if0Default)
	present := isFilePresent(if0File)
	if !present {
		log.Println("Current running configuration missing. Creating a default if0.env file at ~/.if0")
		if _, err := os.Stat(if0Dir); os.IsNotExist(err) {
			log.Println("Directory does not exist, creating dir .if0")
			err = os.Mkdir(if0Dir, os.ModeDir)
			if err != nil {
				log.Fatalln("Error while creating .if0 dir: ", err)
			}
		}
		f, err := os.OpenFile(if0File, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Fatalln("Error while creating a new config file: ", err)
		}
		_, err = f.WriteString("IFO_VERSION=1")
		if err != nil {
			log.Fatalln("Error while writing to the new config file: ", err)
		}
	}
	readConfigFile(if0File)
	allConfig := viper.AllSettings()
	for key, val := range allConfig {
		fmt.Println(key, ":", val)
	}
}

// AddConfigFile replaces the current config file with the provided config file.
// it first checks if the file is already present,
// if present, it creates a backup in the ~if0/.snapshots directory
// and then proceeds to replace the current running config file
func AddConfigFile(srcConfigFile string, zero, merge bool) {
	runningConfigFile := getRunningConfigFile(srcConfigFile, zero)
	// taking a backup of the running configuration if already present
	present := isFilePresent(runningConfigFile)
	if present {
		err := backupToSnapshots(runningConfigFile)
		if err != nil {
			log.Fatalln("Failed to add/update the config file: ", err)
		}
	}
	if merge {
		// TODO: for zero cluster configs, check if the file is already present before merging
		mergeConfigFiles(srcConfigFile, runningConfigFile)
	} else {
		createConfigFile(srcConfigFile, runningConfigFile)
	}
}
