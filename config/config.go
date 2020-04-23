package config

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	IFO_VERSION  = "IF0_VERSION"
	ZERO_VERSION = "ZERO_VERSION"
)

var (
	rootPath, _  = os.UserHomeDir()
	if0Dir       = filepath.Join(rootPath, ".if0")
	envDir       = filepath.Join(if0Dir, ".environments")
	snapshotsDir = filepath.Join(if0Dir, ".snapshots")
	if0Default   = filepath.Join(if0Dir, "if0.env")
)

// SetEnvVariable sets a config variable
func SetEnvVariable(key, value string) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	viper.Set(key, value)
	if value == GetEnvVariable(key) {
		log.Printf("key %s update with %s successful \n", key, value)
	}
}

// GetEnvVariable retrieves the value of a config variable
func GetEnvVariable(key string) string {
	val := viper.Get(key)
	return cast.ToString(val)
}

// PrintCurrentRunningConfig reads the current running if0/env configuration file and prints it
func PrintCurrentRunningConfig() {
	present := isFilePresent(if0Default)
	if !present {
		log.Println("Current running configuration missing. Creating a default if0.env file at ~/.if0")
		if _, err := os.Stat(if0Dir); os.IsNotExist(err) {
			log.Println("Directory does not exist, creating dir .if0")
			err = os.Mkdir(if0Dir, os.ModeDir)
			if err != nil {
				log.Errorln("Error while creating .if0 dir: ", err)
				return
			}
		}
		f, err := os.OpenFile(if0Default, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Errorln("Error while creating a new config file: ", err)
			return
		}
		_, err = f.WriteString("IFO_VERSION=1")
		if err != nil {
			log.Errorln("Error while writing to the new config file: ", err)
			return
		}
	}
	readConfigFile(if0Default)
	allConfig := viper.AllSettings()
	for key, val := range allConfig {
		fmt.Println(key, ":", val)
	}
}

// AddConfigFile replaces the current config file with the provided config file.
// it first checks if the file is already present,
// if present, it creates a backup in the ~if0/.snapshots directory
// and then proceeds to replace the current running config file
func AddConfigFile(srcConfigFile string, zero bool) error {
	runningConfigFile := getRunningConfigFile(srcConfigFile, zero)
	// taking a backup of the running configuration if already present
	present := isFilePresent(runningConfigFile)
	if present {
		err := backupToSnapshots(runningConfigFile)
		if err != nil {
			log.Errorln("Failed to add/update the config file: ", err)
			return err
		}
	}
	err := createConfigFile(srcConfigFile, runningConfigFile)
	if err != nil {
		log.Errorln("Error while creating config file: ", err)
		return err
	}
	return nil
}

// MergeConfigFiles merges configuration at dst with configuration from source.
// it first gets a valid dst file path to be merged with
// if the file is present, the dst file is backed-up in the .snapshots directory
// and then merged with the src config file
func MergeConfigFiles(src, dst string, zero bool) error {
	if src == "" {
		return errors.New("Please provide valid source/destination configuration files for merge.")
	}
	dst = getDstFileForMerge(src, dst, zero)
	srcValid, err := IsConfigFileValid(src, zero)
	if !srcValid {
		log.Println("Please provide a valid configuration file for merge.")
		return err
	}
	if isFilePresent(dst) {
		err = backupToSnapshots(dst)
		if err != nil {
			log.Errorln("Failed to backup the config file: ", err)
			return err
		}
		mergeConfigFiles(src, dst)
	} else {
		log.Errorln("Destination configuration file not found for merge. " +
			"Please provide a valid destination file")
		return err
	}
	return nil
}

// IsConfigFileValid checks if the provided configuration file is valid for the config add/update operation.
// valid if0.env files contain IF0_VERSION key
// valid zero-cluster files contain ZERO_VERSION key
func IsConfigFileValid(configFile string, zero bool) (bool, error) {
	// read IF0_VERSION, ZERO_VERSION
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorln("Error while reading config file: ", err)
	}
	if0Version := viper.IsSet(IFO_VERSION)
	zeroVersion := viper.IsSet(ZERO_VERSION)

	if !if0Version && !zeroVersion {
		return false, errors.New("no valid versions (IF0_VERSION or ZERO_VERSION) found in the config file")
	} else if if0Version && zero {
		return false, errors.New("zero-cluster config update invoked with if0.env config file")
	} else if zeroVersion && !zero {
		return false, errors.New("if0.env update invoked with zero-cluster config file")
	}
	return true, nil
}
