package config

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func mergeConfigFiles(srcConfigFile , runningConfigFile string) {
	readConfigFile(runningConfigFile)
	currentConfigMap := viper.AllSettings()
	readConfigFile(srcConfigFile)
	newConfigMap := viper.AllSettings()
	for k, v := range newConfigMap {
		currentConfigMap[k] = v
	}
	writeToConfigFile(runningConfigFile, currentConfigMap)
}

func writeToConfigFile(runningConfigFile string, currentConfigMap map[string]interface{}) {
	var lines []string
	for key, val := range currentConfigMap {
		lines = append(lines, fmt.Sprintf("%v=%v", strings.ToUpper(key), val))
	}
	s := strings.Join(lines, "\n")
	err := ioutil.WriteFile(runningConfigFile, []byte(s), 0644)
	if err != nil {
		log.Errorln("Error while merging config files: ", err)
		return
	}
}

func getRunningConfigFile(srcConfigFile string, zero bool) string {
	var runningConfigFile string
	if zero {
		// creating a .environments directory to store zero-cluster configurations, if it does not exist.
		if _, err := os.Stat(envDir); os.IsNotExist(err) {
			log.Debugln("Directory does not exist, creating dir for snapshots")
			_ = os.Mkdir(envDir, os.ModeDir)
		}
		// setting configuration file path to update zero-cluster configuration
		log.Println("Updating zero cluster configuration with ", srcConfigFile)
		runningConfigFile = filepath.Join(envDir, filepath.Base(srcConfigFile))
	} else {
		// setting configuration file path to update if0.env configuration
		log.Println("Updating if0.env configuration with ", srcConfigFile)
		runningConfigFile = if0Default
	}
	return runningConfigFile
}

// readConfigFile reads the provided config file
func readConfigFile(configFile string) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorln("Error while reading config file: ", err)
		return
	}
}

// createConfigFile creates a new running config file from the provided config file (src)
func createConfigFile(srcConfigFile, runningConfigFile string) error {
	readConfigFile(srcConfigFile)
	err := viper.WriteConfigAs(runningConfigFile)
	if err != nil {
		log.Errorln("Failed to add/update the config file: ", err)
		return err
	}
	return nil
}

// isFilePresent checks whether the provided config file is already the running config file
// returns false if no, true if yes.
func isFilePresent(fileName string) bool {
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

