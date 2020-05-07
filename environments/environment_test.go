package environments

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/stretchr/testify/assert"
	"if0/common/sync"
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
	err := AddEnv("sample_repo")
	assert.Nil(t, err)
}

func TestSyncEnvNoRepo(t *testing.T) {
	err := SyncEnv("sample-repo")
	assert.EqualError(t, err, "repository not found")
}