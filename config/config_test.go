package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	log.Println("Teardown")
	err := os.RemoveAll("testif0")
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(exitVal)
}

func TestPrintCurrentRunningConfigNoDefaultConfig(t *testing.T) {
	rootPath = "config"
	if0Dir = "testif0"
	_ = os.RemoveAll(if0Dir)
	out := readStdOutPrintCurrentRunningConfig()
	assert.Equal(t, "ifo_version : 1\n", string(out))
}

func TestPrintCurrentRunningConfigWithDefaultConfig(t *testing.T) {
	if0Dir = "testif0"
	out := readStdOutPrintCurrentRunningConfig()
	assert.Equal(t, "ifo_version : 1\n", string(out))
}

func TestAddConfigFileReplace(t *testing.T) {
	if0Dir = "testif0"
	snapshotsDir = filepath.Join(if0Dir, ".snapshots")
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1"), 0644)
	AddConfigFile(testConfig, false, false)
	readConfigFile(filepath.Join(if0Dir, if0Default))
	configMap := viper.AllSettings()
	assert.Equal(t, 1, len(configMap))
	assert.Equal(t, "testval1", configMap["testkey1"])
}

func TestAddConfigFileMerge(t *testing.T) {
	if0Dir = "testif0"
	snapshotsDir = filepath.Join(if0Dir, ".snapshots")
	testConfig := "config2.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey2=testval2"), 0644)
	AddConfigFile(testConfig, false, true)
	readConfigFile(filepath.Join(if0Dir, if0Default))
	configMap := viper.AllSettings()
	assert.Equal(t, 2, len(configMap))
	assert.Equal(t, "testval1", configMap["testkey1"])
	assert.Equal(t, "testval2", configMap["testkey2"])
	_ = os.Remove(testConfig)
}

func TestAddConfigFileEnvironment(t *testing.T) {
	if0Dir = "testif0"
	snapshotsDir = filepath.Join(if0Dir, ".snapshots")
	envDir = filepath.Join(if0Dir, ".environments")
	testConfig := "zero1.env"
	_ = ioutil.WriteFile(testConfig, []byte("zerokey1=zeroval1"), 0644)
	AddConfigFile(testConfig, true, false)
	readConfigFile(filepath.Join(envDir, testConfig))
	configMap := viper.AllSettings()
	assert.Equal(t, 1, len(configMap))
	assert.Equal(t, "zeroval1", configMap["zerokey1"])
	_ = os.Remove(testConfig)
}

//func TestAddConfigFileEnvironmentMerge(t *testing.T) {
//	if0Dir = "testif0"
//	snapshotsDir = filepath.Join(if0Dir, ".snapshots")
//	envDir = filepath.Join(if0Dir, ".environments")
//	testConfig := "zero2.env"
//	_ = ioutil.WriteFile(testConfig, []byte("zerokey2=zeroval2"), 0644)
//	AddConfigFile(testConfig, true, true)
//	readConfigFile(filepath.Join(envDir, testConfig))
//	configMap := viper.AllSettings()
//	assert.Equal(t, 1, len(configMap))
//	assert.Equal(t, "zeroval1", configMap["zerokey1"])
//}

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
	return out
}