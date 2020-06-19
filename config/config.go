package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"if0/common"
	"os"
	"strings"
)

// SetEnvVariable sets a config variable
func SetEnvVariable(key, value string) {
	key = strings.ToUpper(strings.TrimSpace(key))
	value = strings.TrimSpace(value)
	err := os.Setenv(key, value)
	if err != nil {
		fmt.Printf("Error: Setting env variable %s - %s\n", key, err)
	}
	writeToIf0(key, value)
}

// GetEnvVariable retrieves the value of a config variable
func GetEnvVariable(key string) string {
	viper.AutomaticEnv()
	val := viper.Get(key)
	return cast.ToString(val)
}

// PrintCurrentRunningConfig reads the current running if0/env configuration file and prints it
func PrintCurrentRunningConfig() {
	if _, err := os.Stat(common.If0Dir); os.IsNotExist(err) {
		fmt.Println("Directory does not exist, creating dir .if0")
		err = os.Mkdir(common.If0Dir, os.ModeDir)
		if err != nil {
			fmt.Println("Error: Creating .if0 dir - ", err)
			return
		}
	}
	err := writeDefaultIf0Config(common.DefaultEnvFile)
	if err != nil {
		return
	}
	ReadConfigFile(common.If0Default)
	for key, val := range viper.AllSettings() {
		fmt.Println(strings.ToUpper(key)+"="+val.(string))
	}
}

// AddConfigFile replaces the current config file with the provided config file.
// it first checks if the file is already present,
// if present, it creates a backup in the ~if0/.snapshots directory
// and then proceeds to replace the current running config file
func AddConfigFile(srcConfigFile string) error {
	runningConfigFile := common.If0Default
	// taking a backup of the running configuration if already present
	present := isFilePresent(runningConfigFile)
	if present {
		err := backupToSnapshots(runningConfigFile)
		if err != nil {
			fmt.Println("Error: Add/update the config file - ", err)
			return err
		}
	}
	err := createConfigFile(srcConfigFile, runningConfigFile)
	if err != nil {
		fmt.Println("Error: Creating config file - ", err)
		return err
	}
	return nil
}

// MergeConfigFiles merges configuration at dst with configuration from source.
// it first gets a valid dst file path to be merged with
// if the file is present, the dst file is backed-up in the .snapshots directory
// and then merged with the src config file
func MergeConfigFiles(src, dst string) error {
	if src == "" {
		return errors.New("Please provide valid source/destination configuration files for merge.")
	}
	dst = common.If0Default
	srcValid, err := IsConfigFileValid(src)
	if !srcValid {
		fmt.Println("Please provide a valid configuration file for merge.")
		return err
	}
	if isFilePresent(dst) {
		err = backupToSnapshots(dst)
		if err != nil {
			fmt.Println("Error: Config file backup - ", err)
			return err
		}
		mergeConfigFiles(src, dst)
	} else {
		fmt.Println("Destination configuration file not found for merge. " +
			"Please provide a valid destination file")
		return err
	}
	return nil
}

// IsConfigFileValid checks if the provided configuration file is valid for the config add/update operation.
// valid if0.env files contain IF0_VERSION key
func IsConfigFileValid(configFile string) (bool, error) {
	// read IF0_VERSION
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error: Reading config file: ", err)
	}
	if0Version := viper.IsSet(common.IF0_VERSION)

	if !if0Version {
		return false, errors.New("no valid IF0_VERSION found in the config file")
	}
	return true, nil
}
