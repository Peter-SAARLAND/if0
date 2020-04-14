package cmd

import (
	"github.com/stretchr/testify/assert"
	"if0/config"
	"testing"
)

func TestLoadConfigFromParams(t *testing.T) {
	loadConfigFromFlags([]string{"var1=val1","var2=val2"})
	assert.Equal(t, "val1", config.GetEnvVariable("var1"))
	assert.Equal(t, "val2", config.GetEnvVariable("var2"))
}

//func TestLoadConfigFromFileError(t *testing.T) {
//	rescueStdout := os.Stdout
//	r, w, _ := os.Pipe()
//	os.Stdout = w
//	loadConfigFromFile([]string{"config.env"})
//	_ = w.Close()
//	out, _ := ioutil.ReadAll(r)
//	os.Stdout = rescueStdout
//	assert.Contains(t, string(out), `level=error msg="Error while reading config file`)
//}
//
//func TestLoadConfigFromFileInvalid(t *testing.T) {
//	testConfig := "config.env"
//	_ = ioutil.WriteFile(testConfig, []byte("testkey1=testval1\nZERO_VERSION=1"), 0644)
//	rescueStdout := os.Stdout
//	r, w, _ := os.Pipe()
//	os.Stdout = w
//	loadConfigFromFile([]string{"config.env"})
//	_ = w.Close()
//	out, _ := ioutil.ReadAll(r)
//	os.Stdout = rescueStdout
//	fmt.Println("out", string(out))
//	assert.Contains(t, string(out),
//		`Terminating config update:  if0.env update invoked with zero-cluster config file`)
//	_ = os.Remove(testConfig)
//}
