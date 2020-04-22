package config

import (
	"fmt"
	"github.com/djherbis/times"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

// GarbageCollection automatically cleans up backed-up files in the ~/.if0/.snapshots directory
// requires env variables GC_AUTO and GC_PERIOD to be set.
// By default, GC_AUTO=false, GC_PERIOD=30 (days)
func GarbageCollection() {
	gc, gcPeriod := getGcPeriod()
	if gc {
		files, err := ioutil.ReadDir(snapshotsDir)
		if err != nil {
			log.Errorln("error while reading snaphots: ", err)
			return
		}
		for _, f := range files {
			creationTime := getCreationTime(f)
			diff := time.Now().Sub(creationTime).Hours() / 24
			if int(diff) >= gcPeriod {
				_ = os.Remove(filepath.Join(snapshotsDir, f.Name()))
			}
		}
	}
}

// getCreationTime returns the time of creation of a file
// if the time of creation is not available, it returns last modified time
func getCreationTime(f os.FileInfo) time.Time {
	var creationTime time.Time
	t, err := times.Stat(filepath.Join(snapshotsDir, f.Name()))
	if err != nil {
		log.Errorln(err)
	}
	if times.HasBirthTime {
		creationTime = t.BirthTime()
	} else {
		creationTime = f.ModTime()
	}
	return creationTime
}

// getGcPeriod retrieves GC_AUTO and GC_PERIOD from if0.env or environment variables.
// for GC_AUTO = true, automatic garbage collection is triggered,
// files older than GC_PERIOD are deleted
func getGcPeriod() (bool, int) {
	readConfigFile(if0Default)
	gcAutoStr := GetEnvVariable("GC_AUTO")
	gcPeriodStr := GetEnvVariable("GC_PERIOD")

	gcAuto := parseGcAuto(gcAutoStr)
	// if GC_PERIOD is not set, setting it to default value of 30 days
	var gcPeriod int
	if gcAuto {
		if gcPeriodStr == "" {
			gcPeriod = 30
		} else {
			gcPeriod, _ = strconv.Atoi(gcPeriodStr)
		}
		return true, gcPeriod
	}
	return false, 0
}

// parseGcAuto parses GC_AUTO variable for values:
// t, true, True, TRUE, yes, Yes, YES, 1
// f, false, False, FALSE, no, No, NO, 0
func parseGcAuto(gcAutoStr string) bool {
	var gcAuto bool
	if strings.ToLower(gcAutoStr) == "yes" {
		gcAuto = true
	} else if strings.ToLower(gcAutoStr) == "no" {
		gcAuto = false
	} else {
		gcAuto, _ = strconv.ParseBool(gcAutoStr)
	}
	return gcAuto
}
