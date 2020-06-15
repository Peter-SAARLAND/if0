package environments

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/stretchr/testify/assert"
	"if0/common"
	"if0/common/sync"
	"io/ioutil"
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

func TestListEnv(t *testing.T) {
	common.EnvDir = "testdata"
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	ListEnv()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = stdout
	localRepo := filepath.Join("testdata", "test-env-1")
	assert.Contains(t, string(out), "Local repository: "+localRepo)
	assert.Contains(t, string(out), "Repository URL:  remote repository does not exist\n")
}

func TestInspectEnv(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	InspectEnv("testdata/test-env-1")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = stdout
	assert.Contains(t, string(out), "IF0_ENVIRONMENT=test-repo-1\n")
}