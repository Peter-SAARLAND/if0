package environments

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/stretchr/testify/assert"
	"if0/common"
	"if0/common/sync"
	"if0/config"
	"os"
	"path/filepath"
	"testing"
)

func TestAddEnvAuthError(t *testing.T) {
	getAuth = func(authObj sync.AuthOps, remoteStorage string) (transport.AuthMethod, error) {
		return nil, errors.New("test-auth-error")
	}
	err := AddEnv("sample_repo")
	assert.EqualError(t, err, "test-auth-error")
}

func TestAddEnvCloneError(t *testing.T) {
	getAuth = func(authObj sync.AuthOps, remoteStorage string) (transport.AuthMethod, error) {
		return nil, nil
	}
	clone = func(repoUrl string, auth transport.AuthMethod) (*git.Repository, error) {
		return nil, errors.New("test-clone-error")
	}
	err := AddEnv("sample_repo")
	assert.EqualError(t, err, "test-clone-error")
}

func TestAddEnvClone(t *testing.T) {
	getAuth = func(authObj sync.AuthOps, remoteStorage string) (transport.AuthMethod, error) {
		return nil, nil
	}
	clone = func(repoUrl string, auth transport.AuthMethod) (*git.Repository, error) {
		r := &git.Repository{}
		return r, nil
	}
	pushEnvInitChanges = func(r *git.Repository, auth transport.AuthMethod) error {
		return nil
	}
	err := AddEnv("sample_repo")
	assert.Nil(t, err)
}

func TestSyncEnvNoRepo(t *testing.T) {
	err := SyncEnv("sample-repo")
	assert.EqualError(t, err, "repository not found")
}

func TestSyncEnvError(t *testing.T) {
	common.EnvDir = "testdata"
	_ = os.Mkdir("testdata", 0644)
	_ = os.Mkdir(filepath.Join("testdata", "sample-repo"), 0644)
	repoSync = func(syncObj sync.SyncOps, repo string, if0Repo bool) error {
		return errors.New("test-repo-sync-error")
	}
	err := SyncEnv(filepath.Join("testdata", "sample-repo"))
	assert.EqualError(t, err, "test-repo-sync-error")
}

func TestSyncEnv(t *testing.T) {
	common.EnvDir = "testdata"
	_ = os.Mkdir("testdata", 0644)
	_ = os.Mkdir(filepath.Join("testdata", "sample-repo"), 0644)
	repoSync = func(syncObj sync.SyncOps, repo string, if0Repo bool) error {
		return nil
	}
	err := SyncEnv(filepath.Join("testdata", "sample-repo"))
	assert.Nil(t, err)
}

func TestLoadEnvNoFiles(t *testing.T) {
	common.EnvDir = filepath.Join("testdata", "sample-repo")
	os.Remove(filepath.Join("testdata", "sample-repo", "if0.env"))
	err := loadEnv(common.EnvDir)
	assert.EqualError(t, err, "no .env files found")
}

func TestLoadEnv(t *testing.T) {
	common.EnvDir = filepath.Join("testdata", "sample-repo")
	f, _ := os.OpenFile(filepath.Join("testdata", "sample-repo", "if0.env"), os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	_, _ = f.Write([]byte("IF0_ENVIRONMENT=sample-repo"))
	_ = loadEnv(common.EnvDir)
	assert.Equal(t, "sample-repo", config.GetEnvVariable("IF0_ENVIRONMENT"))
}

func TestEnvInit(t *testing.T) {
	common.EnvDir = "testdata"
	common.If0Default = filepath.Join("testdata", "if0.env")
	pushEnvInitChanges = func(r *git.Repository, auth transport.AuthMethod) error {
		return nil
	}
	err := envInit(nil, nil, "sample-repo")
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