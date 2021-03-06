package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"if0/common"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPrintCurrentRunningConfigNoDefaultConfig(t *testing.T) {
	common.RootPath = "config"
	common.If0Dir = "testif0"
	common.If0Default = filepath.Join(common.If0Dir, "if0.env")
	common.DefaultEnvFile = filepath.Join("defenv", "defaultIf0.env")
	_ = os.RemoveAll(common.If0Dir)
	PrintCurrentRunningConfig()
	ReadConfigFile(common.If0Default)
	assert.Equal(t, "1", GetEnvVariable("IF0_VERSION"))
	assert.Equal(t, "https://gitlab.com/peter.saarland/shipmate/-/raw/master/shipmate.gitlab-ci.yml",
		GetEnvVariable("SHIPMATE_WORKFLOW_URL"))
}

func TestPrintCurrentRunningConfigWithDefaultConfig(t *testing.T) {
	common.If0Dir = "testif0"
	common.If0Default = filepath.Join(common.If0Dir, "if0.env")
	common.DefaultEnvFile = filepath.Join("defenv", "defaultIf0.env")
	PrintCurrentRunningConfig()
	ReadConfigFile(common.If0Default)
	assert.Equal(t, "1", GetEnvVariable("IF0_VERSION"))
	assert.Equal(t, "https://gitlab.com/peter.saarland/shipmate/-/raw/master/shipmate.gitlab-ci.yml",
		GetEnvVariable("SHIPMATE_WORKFLOW_URL"))
}

func TestAddConfigFileReplace(t *testing.T) {
	common.If0Dir = "testif0"
	common.If0Default = filepath.Join(common.If0Dir, "if0.env")
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1"), 0644)
	AddConfigFile(testConfig)
	ReadConfigFile(common.If0Default)
	configMap := viper.AllSettings()
	assert.Equal(t, 1, len(configMap))
	assert.Equal(t, "testval1", configMap["testkey1"])
}

func TestMergeConfigFiles(t *testing.T) {
	common.If0Dir = "testif0"
	common.If0Default = filepath.Join(common.If0Dir, "if0.env")
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	testConfig := "config2.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey2=testval2\nIF0_VERSION=1"), 0644)
	_ = MergeConfigFiles(testConfig, common.If0Default)
	ReadConfigFile(common.If0Default)
	configMap := viper.AllSettings()
	assert.Equal(t, 3, len(configMap))
	assert.Equal(t, "testval1", configMap["testkey1"])
	assert.Equal(t, "testval2", configMap["testkey2"])
	_ = os.Remove(testConfig)
}

func TestAddConfigFileEnvironment(t *testing.T) {
	common.If0Dir = "testif0"
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	common.EnvDir = filepath.Join(common.If0Dir, ".environments")
	testConfig := "zero1.env"
	_ = ioutil.WriteFile(testConfig, []byte("zerokey1=zeroval1"), 0644)
	AddConfigFile(testConfig)
	ReadConfigFile(filepath.Join(common.EnvDir, testConfig))
	configMap := viper.AllSettings()
	assert.Equal(t, 1, len(configMap))
	assert.Equal(t, "zeroval1", configMap["zerokey1"])
	_ = os.Remove(testConfig)
}

func TestMergeConfigFilesInvalid(t *testing.T) {
	common.If0Dir = "testif0"
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	common.EnvDir = filepath.Join(common.If0Dir, ".environments")
	common.If0Default = filepath.Join(common.If0Dir, "if0.env")
	err := MergeConfigFiles("abc.env", "")
	assert.Error(t, err)
}


func TestSetEnvVariable(t *testing.T) {
	SetEnvVariable("test", "val")
	val := GetEnvVariable("test")
	fmt.Println(viper.AllSettings())
	assert.Equal(t, "val", val)
}

func TestIsConfigFileValidTrue(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nIF0_VERSION=1"), 0644)
	isValid, _ := IsConfigFileValid(testConfig)
	assert.Equal(t, isValid, true)
	_ = os.Remove(testConfig)
}

func TestIsConfigFileValidFalse(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nabc=def"), 0644)
	isValid, _ := IsConfigFileValid(testConfig)
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestIf0WithZeroConfig(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nZERO_VERSION=1"), 0644)
	isValid, err := IsConfigFileValid(testConfig)
	assert.Error(t, err)
	assert.Equal(t, "no valid IF0_VERSION found in the config file", err.Error())
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestNoValidConfig(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\n"), 0644)
	isValid, err := IsConfigFileValid(testConfig)
	assert.Error(t, err)
	assert.Equal(t, "no valid IF0_VERSION found in the config file", err.Error())
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)

}

func TestMergeConfigFilesNoSrc(t *testing.T) {
	err := MergeConfigFiles("", "")
	assert.Equal(t, "Please provide valid source/destination configuration files for merge.", err.Error())
}

func TestGarbageCollectionNo(t *testing.T) {
	common.If0Dir = "testif0"
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	SetEnvVariable("GC_AUTO", "No")
	SetEnvVariable("GC_PERIOD", "0")
	GarbageCollection()
	f, _ := ioutil.ReadDir(common.SnapshotsDir)
	assert.NotEqual(t, 0, len(f))
}

func TestGarbageCollection(t *testing.T) {
	common.If0Dir = "testif0"
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	SetEnvVariable("GC_AUTO", "Yes")
	SetEnvVariable("GC_PERIOD", "0")
	GarbageCollection()
	f, _ := ioutil.ReadDir(common.SnapshotsDir)
	assert.Equal(t, 0, len(f))
	_ = os.RemoveAll("testif0")
}

func TestParseGcAutoStr(t *testing.T) {
	assert.Equal(t, true, parseGcAuto("1"))
	assert.Equal(t, true, parseGcAuto("true"))
	assert.Equal(t, true, parseGcAuto("t"))
	assert.Equal(t, true, parseGcAuto("TRUE"))
	assert.Equal(t, true, parseGcAuto("True"))
	assert.Equal(t, true, parseGcAuto("yes"))
	assert.Equal(t, true, parseGcAuto("YES"))
	assert.Equal(t, false, parseGcAuto("0"))
	assert.Equal(t, false, parseGcAuto("false"))
	assert.Equal(t, false, parseGcAuto("f"))
	assert.Equal(t, false, parseGcAuto("False"))
	assert.Equal(t, false, parseGcAuto("FALSE"))
	assert.Equal(t, false, parseGcAuto("no"))
	assert.Equal(t, false, parseGcAuto("NO"))
}

func TestWriteDefaultIf0ConfigNone(t *testing.T) {
	common.If0Default = filepath.Join("testdata", "if0.env")
	defFile := filepath.Join("testdata", "testDefEnv.env")
	ioutil.WriteFile(defFile, []byte("IF0_VERSION=1\nTESTIF0=YES\n"), 0644)
	os.Remove(common.If0Default)
	err := writeDefaultIf0Config(defFile)
	assert.Nil(t, err)
	assert.FileExists(t, common.If0Default)
	defBytes, _ := ioutil.ReadFile(defFile)
	newIf0Bytes, _ := ioutil.ReadFile(common.If0Default)
	assert.Equal(t, defBytes, newIf0Bytes)
}

func TestWriteDefaultIf0ConfigAppend(t *testing.T) {
	common.If0Default = filepath.Join("testdata", "if0.env")
	defFile := filepath.Join("testdata", "testDefEnv.env")
	ioutil.WriteFile(defFile, []byte("IF0_VERSION=1\nTESTIF0=YES\nAPPEND=VAL\n"), 0644)
	err := writeDefaultIf0Config(defFile)
	assert.Nil(t, err)
	assert.FileExists(t, common.If0Default)
	newIf0Bytes, _ := ioutil.ReadFile(common.If0Default)
	assert.Contains(t, string(newIf0Bytes), "APPEND=VAL")
	ioutil.WriteFile(defFile, []byte("IF0_VERSION=1\nTESTIF0=YES\n"), 0644)
	os.Remove(common.If0Default)
}
