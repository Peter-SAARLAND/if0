package environments

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/stretchr/testify/assert"
	"if0/common"
	"os"
	"path/filepath"
	"testing"
)

func TestEnvInit(t *testing.T) {
	common.EnvDir = "testdata"
	common.If0Default = filepath.Join("testdata", "if0.env")
	pushEnvInitChanges = func(r *git.Repository, auth transport.AuthMethod) error {
		return nil
	}
	err := envInit(filepath.Join("testdata", "sample-repo"))
	assert.Nil(t, err)
	assert.DirExists(t, filepath.Join("testdata", "sample-repo", ".ssh"))
	assert.FileExists(t, filepath.Join("testdata", "sample-repo", "zero.env"))
	assert.FileExists(t, filepath.Join("testdata", "sample-repo", ".gitlab-ci.yml"))
	os.RemoveAll(filepath.Join("testdata", "sample-repo", ".ssh"))
	os.Remove(filepath.Join("testdata", "sample-repo", "zero.env"))
	os.Remove(filepath.Join("testdata", "sample-repo", ".gitlab-ci.yml"))
}

func TestGetShipmateUrl(t *testing.T) {
	common.If0Default = filepath.Join("testdata", "if0.env")
	assert.Equal(t, "https://gitlab.com/peter.saarland/shipmate/-/raw/master/shipmate.gitlab-ci.yml", getShipmateUrl())
}

func TestCreateNestedDirPathWithUrl(t *testing.T) {
	dirPath := createNestedDirPath("test-env-1", "git@gitlab.com:grvijayan/test-env-1.git")
	expected := filepath.Join(common.EnvDir, "gitlab.com", "grvijayan", "test-env-1")
	assert.Equal(t, expected, dirPath)
}

func TestCreateNestedDirPathWithNoUrl(t *testing.T) {
	dirPath := createNestedDirPath("test-env-1", "")
	expected := filepath.Join(common.EnvDir, "test-env-1")
	assert.Equal(t, expected, dirPath)
}

func TestAddLocalEnv(t *testing.T) {
	envDir := filepath.Join("testdata", "gitlab.com", "peter-saarland", "test-env-2")
	addLocalEnv(envDir)
	assert.DirExists(t, filepath.Join("testdata", "gitlab.com"))
	assert.DirExists(t, filepath.Join("testdata", "gitlab.com", "peter-saarland"))
	assert.DirExists(t, filepath.Join("testdata", "gitlab.com", "peter-saarland", "test-env-2"))
	os.RemoveAll(filepath.Join("testdata", "gitlab.com"))
}

func TestReadEnvNoFiles(t *testing.T) {
	common.EnvDir = filepath.Join("testdata", "sample-repo")
	os.Remove(filepath.Join("testdata", "sample-repo", "if0.env"))
	os.Remove(filepath.Join("testdata", "sample-repo", "logo.png"))
	allConfig, err := readAllEnvFiles(common.EnvDir)
	assert.EqualError(t, err, "no .env files found")
	assert.Nil(t, allConfig)
}

func TestRealAllEnv(t *testing.T) {
	common.EnvDir = filepath.Join("testdata", "sample-repo")
	f, _ := os.OpenFile(filepath.Join("testdata", "sample-repo", "if0.env"), os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	_, _ = f.Write([]byte("IF0_ENVIRONMENT=sample-repo"))
	allConfig, _ := readAllEnvFiles(common.EnvDir)
	assert.Equal(t, allConfig["if0_environment"], "sample-repo")
}