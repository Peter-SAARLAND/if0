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
	_ = os.RemoveAll(common.If0Dir)
	out := readStdOutPrintCurrentRunningConfig()
	assert.Contains(t, string(out), "ifo_version : 1\n")
}

func TestPrintCurrentRunningConfigWithDefaultConfig(t *testing.T) {
	common.If0Dir = "testif0"
	common.If0Default = filepath.Join(common.If0Dir, "if0.env")
	out := readStdOutPrintCurrentRunningConfig()
	assert.Equal(t, "ifo_version : 1\n", string(out))
}

func TestAddConfigFileReplace(t *testing.T) {
	common.If0Dir = "testif0"
	common.If0Default = filepath.Join(common.If0Dir, "if0.env")
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1"), 0644)
	AddConfigFile(testConfig, false)
	readConfigFile(common.If0Default)
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
	_ = MergeConfigFiles(testConfig, common.If0Default, false)
	readConfigFile(common.If0Default)
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
	AddConfigFile(testConfig, true)
	readConfigFile(filepath.Join(common.EnvDir, testConfig))
	configMap := viper.AllSettings()
	assert.Equal(t, 1, len(configMap))
	assert.Equal(t, "zeroval1", configMap["zerokey1"])
	_ = os.Remove(testConfig)
}

func TestMergeFilesEnvironment(t *testing.T) {
	common.If0Dir = "testif0"
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	common.EnvDir = filepath.Join(common.If0Dir, ".environments")
	testConfig := "zero2.env"
	_ = ioutil.WriteFile(testConfig, []byte("zerokey2=zeroval2\nZERO_VERSION=2"), 0644)
	_ = MergeConfigFiles(testConfig, "zero1.env", true)
	readConfigFile(filepath.Join(common.EnvDir, "zero1.env"))
	configMap := viper.AllSettings()
	assert.Equal(t, 3, len(configMap))
	assert.Equal(t, "zeroval1", configMap["zerokey1"])
	assert.Equal(t, "zeroval2", configMap["zerokey2"])
	_ = os.Remove(testConfig)
}

func TestMergeConfigFilesInvalid(t *testing.T) {
	common.If0Dir = "testif0"
	common.SnapshotsDir = filepath.Join(common.If0Dir, ".snapshots")
	common.EnvDir = filepath.Join(common.If0Dir, ".environments")
	common.If0Default = filepath.Join(common.If0Dir, "if0.env")
	err := MergeConfigFiles("abc.env", "", false)
	assert.Error(t, err)
}


func TestSetEnvVariable(t *testing.T) {
	SetEnvVariable("test", "val")
	val := GetEnvVariable("test")
	fmt.Println(viper.AllSettings())
	assert.Equal(t, "val", val)
}

func readStdOutPrintCurrentRunningConfig() []byte {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	PrintCurrentRunningConfig()
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	_ = r.Close()
	return out
}

func TestIsConfigFileValidTrue(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nIF0_VERSION=1"), 0644)
	isValid, _ := IsConfigFileValid(testConfig, false)
	assert.Equal(t, isValid, true)
	_ = os.Remove(testConfig)
}

func TestIsConfigFileValidFalse(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nabc=def"), 0644)
	isValid, _ := IsConfigFileValid(testConfig, false)
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestIsZeroConfigFileValidTrue(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nZERO_VERSION=1"), 0644)
	isValid, _ := IsConfigFileValid(testConfig, true)
	assert.Equal(t, isValid, true)
	_ = os.Remove(testConfig)
}

func TestIsZeroConfigFileValidFalse(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\n11ZERO_VERSION=1"), 0644)
	isValid, _ := IsConfigFileValid(testConfig, true)
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestIf0WithZeroConfig(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nZERO_VERSION=1"), 0644)
	isValid, err := IsConfigFileValid(testConfig, false)
	assert.Error(t, err)
	assert.Equal(t, "if0.env update invoked with zero-cluster config file", err.Error())
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestZeroWithIf0Config(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nIF0_VERSION=1"), 0644)
	isValid, err := IsConfigFileValid(testConfig, true)
	assert.Error(t, err)
	assert.Equal(t, "zero-cluster config update invoked with if0.env config file", err.Error())
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestNoValidConfig(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\n"), 0644)
	isValid, err := IsConfigFileValid(testConfig, false)
	assert.Error(t, err)
	assert.Equal(t, "no valid versions (IF0_VERSION or ZERO_VERSION) found in the config file", err.Error())
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)

}

func TestMergeConfigFilesNoSrc(t *testing.T) {
	err := MergeConfigFiles("", "", false)
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
