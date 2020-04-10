package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"if0/config"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfigFromParams(t *testing.T) {
	loadConfigFromFlags([]string{"var1=val1","var2=val2"})
	assert.Equal(t, "val1", config.GetEnvVariable("var1"))
	assert.Equal(t, "val2", config.GetEnvVariable("var2"))
}

func TestLoadConfigFromFileError(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	loadConfigFromFile([]string{"config.env"})
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Contains(t, out, `level=fatal msg="Error while reading config file`)
}

func TestLoadConfigFromFileInvalid(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nZERO_VERSION=1"), 0644)
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	loadConfigFromFile([]string{"config.env"})
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	assert.Contains(t, out, `level=fatal msg="Terminating config update`)
	_ = os.Remove(testConfig)
}

func TestIsConfigFileValidTrue(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nIF0_VERSION=1"), 0644)
	isValid, _ := isConfigFileValid(testConfig)
	assert.Equal(t, isValid, true)
	_ = os.Remove(testConfig)
}

func TestIsConfigFileValidFalse(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nabc=def"), 0644)
	isValid, _ := isConfigFileValid(testConfig)
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestIsZeroConfigFileValidTrue(t *testing.T) {
	testConfig := "config.env"
	zero = true
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nZERO_VERSION=1"), 0644)
	isValid, _ := isConfigFileValid(testConfig)
	assert.Equal(t, isValid, true)
	_ = os.Remove(testConfig)
}

func TestIsZeroConfigFileValidFalse(t *testing.T) {
	testConfig := "config.env"
	zero = true
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\n11ZERO_VERSION=1"), 0644)
	isValid, _ := isConfigFileValid(testConfig)
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestIf0WithZeroConfig(t *testing.T) {
	testConfig := "config.env"
	zero = false
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nZERO_VERSION=1"), 0644)
	isValid, err := isConfigFileValid(testConfig)
	fmt.Println(err)
	assert.Error(t, err)
	assert.Equal(t, "if0.env update invoked with zero-cluster config file", err.Error())
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestZeroWithIf0Config(t *testing.T) {
	testConfig := "config.env"
	zero = true
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nIF0_VERSION=1"), 0644)
	isValid, err := isConfigFileValid(testConfig)
	assert.Error(t, err)
	assert.Equal(t, "zero-cluster config update invoked with if0.env config file", err.Error())
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}

func TestNoValidConfig(t *testing.T) {
	testConfig := "config.env"
	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\n"), 0644)
	isValid, err := isConfigFileValid(testConfig)
	assert.Error(t, err)
	assert.Equal(t, "no valid versions (IF0_VERSION or ZERO_VERSION) found in the config file", err.Error())
	assert.Equal(t, isValid, false)
	_ = os.Remove(testConfig)
}
