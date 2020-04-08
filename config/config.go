package config

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
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
func AddConfigFile(srcConfigFile string, zero bool) {
	srcConfigFileName := filepath.Base(srcConfigFile)
	var runningConfigFile string
	if zero {
		// creating a .environments directory to store zero-cluster configurations, if it does not exist.
		if _, err := os.Stat(envDir); os.IsNotExist(err) {
			log.Debugln("Directory does not exist, creating dir for snapshots")
			_ = os.Mkdir(envDir, os.ModeDir)
		}
		// setting configuration file path to update zero-cluster configuration
		log.Println("Updating zero cluster configuration with ", srcConfigFile)
		runningConfigFile = filepath.Join(envDir, srcConfigFileName)
	} else {
		// setting configuration file path to update if0.env configuration
		log.Println("Updating if0.env configuration with ", srcConfigFile)
		runningConfigFile = filepath.Join(rootPath, if0Default)
	}

	// taking a backup of the running configuration if already present
	present := isFilePresent(runningConfigFile, zero)
	if present {
		err := backupToSnapshots(runningConfigFile)
		if err != nil {
			log.Fatalln("Failed to add/update the config file: ", err)
		}
	}
	createConfigFile(srcConfigFile, runningConfigFile)
}

// readConfigFile reads the provided config file
func readConfigFile(configFile string) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Error while reading config file: ", err)
	}
}

// createConfigFile creates a new running config file from the provided config file (src)
func createConfigFile(srcConfigFile, runningConfigFile string) {
	readConfigFile(srcConfigFile)
	err := viper.WriteConfigAs(runningConfigFile)
	if err != nil {
		log.Fatalln("Failed to add/update the config file: ", err)
	}
}

// isFilePresent checks whether the provided config file is already the running config file
// returns false if no, true if yes.
func isFilePresent(fileName string, zero bool) bool {
	//filePath := rootPath + string(os.PathSeparator) + fileName
	//zeroConfigFilePath := filepath.Join(rootPath, ".environments")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

// backupToSnapshots takes a backup of the current running-config.env file, say if0.env
// it creates a copy of if0.env file, with timestamp attached to the filename
// and stores it in the ~if0/.snapshots directory
// example: if0-02042020_170240.env
func backupToSnapshots(fileName string) error {
	if _, err := os.Stat(snapshotsDir); os.IsNotExist(err) {
		log.Debugln("Directory does not exist, creating dir for snapshots")
		_ = os.Mkdir(snapshotsDir, os.ModeDir)
	}
	timestamp := string(time.Now().Format("02012006_150405"))
	bkpFile := filepath.Join(snapshotsDir, strings.Split(filepath.Base(fileName), ".")[0]+"-"+timestamp+".env")
	readConfigFile(fileName)
	err := viper.WriteConfigAs(bkpFile)
	if err != nil {
		log.Errorln("Error while writing to backup file: ", err)
		return errors.New("backup of previous config failed")
	}
	return nil
}
