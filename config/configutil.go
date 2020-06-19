package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"if0/common"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// getDstFileForMerge returns dst filepath based on zero flag.
// for zero = true, if dst is not provided, dst file name is set to src file name,
// and looked-up in the .environments directory
// if zero is set to false, the src and dst files are the same - if0.env
func getDstFileForMerge(src string, dst string, zero bool) string {
	// setting dst path for zero configuration files
	if zero {
		if dst == "" {
			dst = filepath.Join(common.EnvDir, filepath.Base(src))
		} else {
			dst = filepath.Join(common.EnvDir, filepath.Base(dst))
		}
	} else {
		dst = common.If0Default
	}
	return dst
}

// mergeConfigFiles combines configuration from source .env file with configuration in the destination .env file
// For config keys that are already present, the values are updated from source .env file
func mergeConfigFiles(srcConfigFile, dstConfigFile string) {
	ReadConfigFile(dstConfigFile)
	currentConfigMap := viper.AllSettings()
	ReadConfigFile(srcConfigFile)
	newConfigMap := viper.AllSettings()
	for k, v := range newConfigMap {
		currentConfigMap[k] = v
	}
	writeToConfigFile(dstConfigFile, currentConfigMap)
}

// writeToConfigFile is used to merge config files.
// key-value pairs from currentConfigMap written to the current running config file
func writeToConfigFile(runningConfigFile string, currentConfigMap map[string]interface{}) {
	var lines []string
	for key, val := range currentConfigMap {
		lines = append(lines, fmt.Sprintf("%v=%v", strings.ToUpper(key), val))
	}
	s := strings.Join(lines, "\n")
	err := ioutil.WriteFile(runningConfigFile, []byte(s), 0644)
	if err != nil {
		fmt.Println("Error: Merging config files - ", err)
		return
	}
}

// getRunningConfigFile returns the configuration file to be backed-up and updated
// if the zero flag is set to true, configuration from .environments directory is set as the running config file
// if the zero flag is set to false, if0.env is set as the running configuration file
func getRunningConfigFile(srcConfigFile string, zero bool) string {
	var runningConfigFile string
	if zero {
		// creating a .environments directory to store zero-cluster configurations, if it does not exist.
		if _, err := os.Stat(common.EnvDir); os.IsNotExist(err) {
			fmt.Println("Directory does not exist, creating dir for snapshots")
			_ = os.Mkdir(common.EnvDir, os.ModeDir)
		}
		// setting configuration file path to update zero-cluster configuration
		fmt.Println("Updating zero cluster configuration with ", srcConfigFile)
		runningConfigFile = filepath.Join(common.EnvDir, filepath.Base(srcConfigFile))
	} else {
		// setting configuration file path to update if0.env configuration
		fmt.Println("Updating if0.env configuration with ", srcConfigFile)
		runningConfigFile = common.If0Default
	}
	return runningConfigFile
}

// ReadConfigFile reads the provided config file
func ReadConfigFile(configFile string) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error: while reading config file - ", err)
		return
	}
}

// createConfigFile creates a new running config file from the provided config file (src)
func createConfigFile(srcConfigFile, runningConfigFile string) error {
	ReadConfigFile(srcConfigFile)
	err := viper.WriteConfigAs(runningConfigFile)
	if err != nil {
		fmt.Println("Error: Failed to add/update the config file - ", err)
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
	if _, err := os.Stat(common.SnapshotsDir); os.IsNotExist(err) {
		fmt.Println("Directory does not exist, creating dir for snapshots")
		_ = os.Mkdir(common.SnapshotsDir, os.ModeDir)
	}
	timestamp := string(time.Now().Format("02012006_150405"))
	bkpFile := filepath.Join(common.SnapshotsDir, strings.Split(filepath.Base(fileName), ".")[0]+"-"+timestamp+".env")
	ReadConfigFile(fileName)
	err := viper.WriteConfigAs(bkpFile)
	if err != nil {
		fmt.Println("Error: while writing to backup file - ", err)
		return errors.New("backup of previous config failed")
	}
	return nil
}

// writeDefaultIf0Config creates an if0.env file if not present at ~/.if0/
// and copies the contents of defenv/defaultIf0.env to if0.env.
// If if0.env is present at ~/.if0/, it appends the new contents from defenv/defaultIf0.env to if0.env
// This requires the user to run 'if0 config'
func writeDefaultIf0Config(defaultEnvFile string) error {
	defEnvBytes, err := ioutil.ReadFile(defaultEnvFile)
	if err != nil {
		fmt.Println("Error: Reading default .env file - ", err)
		return err
	}

	if _, err := os.Stat(common.If0Default); os.IsNotExist(err) {
		fmt.Println("if0.env does not exist, creating ", common.If0Default)
		err = ioutil.WriteFile(common.If0Default, defEnvBytes, 0644)
		if err != nil {
			fmt.Println("Error: Writing to if0.env file - ", err)
			return err
		}
	} else {
		mergeConfigFiles(defaultEnvFile, common.If0Default)
	}
	return nil
}

func writeToIf0(key, value string) {
	s := key + "=" + value + "\n"
	err := ioutil.WriteFile(common.If0Default, []byte(s), 0644)
	if err != nil {
		fmt.Printf("Error: Setting %s in if0.env -%s\n", key, err)
		return
	}
}
