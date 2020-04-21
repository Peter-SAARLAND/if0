package config

import (
	"github.com/djherbis/times"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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
			creationTime := getCreationTome(f)
			diff := time.Now().Sub(creationTime).Hours() / 24
			if int(diff) >= gcPeriod {
				_ = os.Remove(filepath.Join(snapshotsDir, f.Name()))
			}
		}
	}
}

func getCreationTome(f os.FileInfo) time.Time {
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
