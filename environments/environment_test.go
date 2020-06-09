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
	err := AddEnv([]string{"add", "sample_repo", "repo-url"})
	assert.EqualError(t, err, "test-auth-error")
}

func TestAddEnvClone(t *testing.T) {
	getAuth = func(authObj sync.AuthOps, remoteStorage string) (transport.AuthMethod, error) {
		return nil, nil
	}

	pushEnvInitChanges = func(r *git.Repository, auth transport.AuthMethod) error {
		return nil
	}
	err := AddEnv([]string{"add", "sample_repo", "repo-url"})
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
	repoSync = func(syncObj sync.SyncOps, repo string, dir string) error {
		return errors.New("test-repo-sync-error")
	}
	err := SyncEnv(filepath.Join("testdata", "sample-repo"))
	assert.EqualError(t, err, "test-repo-sync-error")
}

func TestSyncEnv(t *testing.T) {
	common.EnvDir = "testdata"
	_ = os.Mkdir("testdata", 0644)
	_ = os.Mkdir(filepath.Join("testdata", "sample-repo"), 0644)
	repoSync = func(syncObj sync.SyncOps, repo string, dir string) error {
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

