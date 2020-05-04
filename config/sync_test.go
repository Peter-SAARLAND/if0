package config

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"if0/common"
	"os"
	"path/filepath"
	"testing"
)

type mockSync struct {
	mock.Mock
}

func (m *mockSync) Clone(repoUrl string, auth transport.AuthMethod) (*git.Repository, error) {
	panic("implement me")
}

func (m *mockSync) GitInit(localRepoPath string) (*git.Repository, error) {
	args := m.Called()
	return args.Get(0).(*git.Repository), args.Error(1)
}

func (m *mockSync) AddRemote(remoteStorage string, r *git.Repository) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockSync) Open(if0Dir string) (*git.Repository, error) {
	args := m.Called()
	return nil, args.Error(1)
}

func (m *mockSync) Pull(remoteStorage string, r *git.Repository, pullOptions *git.PullOptions) (*git.Worktree, error) {
	args := m.Called()
	return args.Get(0).(*git.Worktree), args.Error(1)
}

func (m *mockSync) Status(w *git.Worktree) (git.Status, error) {
	args := m.Called()
	return args.Get(0).(git.Status), args.Error(1)
}

func (m *mockSync) AddFile(w *git.Worktree, file string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockSync) Commit(w *git.Worktree) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockSync) Push(auth transport.AuthMethod, r *git.Repository) error {
	args := m.Called()
	return args.Error(0)
}

func TestGitSyncAuthError(t *testing.T) {
	testSyncObj := new(mockSync)
	common.GetSyncAuth = func(authObj common.AuthOps, remoteStorage string) (transport.AuthMethod, error) {
		return nil, errors.New("test-auth-error")
	}
	err := gitSync(testSyncObj, "http://sample-storage")
	assert.EqualError(t, err, "test-auth-error")
}

func TestGitSyncInitError(t *testing.T) {
	testSyncObj := new(mockSync)
	common.If0Dir = "config"
	common.GetSyncAuth = func(authObj common.AuthOps, remoteStorage string) (transport.AuthMethod, error) {
		return nil, nil
	}
	testSyncObj.On("GitInit").Return(&git.Repository{}, errors.New("test-init-error"))
	err := gitSync(testSyncObj, "http://sample-storage")
	assert.EqualError(t, err, "test-init-error")
}

func TestGitInit(t *testing.T) {
	common.GetSyncAuth = func(authObj common.AuthOps, remoteStorage string) (transport.AuthMethod, error) {
		auth := &http.BasicAuth{Username: "test-user", Password: "test-password"}
		return auth, nil
	}
	common.If0Dir = "config"
	var testSyncObj common.Sync
	err := gitSync(&testSyncObj, "http://sample-storage")
	gitDir := filepath.Join("config", ".git")
	assert.DirExists(t, gitDir)
	assert.Nil(t, err)
	_ = os.RemoveAll(common.If0Dir)
}

func TestGitSyncRemoteError(t *testing.T) {
	testSyncObj := new(mockSync)
	common.If0Dir = "config"
	common.GetSyncAuth = func(authObj common.AuthOps, remoteStorage string) (transport.AuthMethod, error) {
		return nil, nil
	}
	testSyncObj.On("GitInit").Return(&git.Repository{}, nil)
	testSyncObj.On("AddRemote").Return(errors.New("test-remote-error"))
	err := gitSync(testSyncObj, "http://sample-storage")
	assert.EqualError(t, err, "test-remote-error")
}

func TestRepoSyncNoRemoteStorage(t *testing.T) {
	SetEnvVariable("REMOTE_STORAGE", "")
	err := RepoSync()
	assert.EqualError(t, err, "REMOTE_STORAGE is not set.")
}

func TestRepoSyncError(t *testing.T) {
	SetEnvVariable("REMOTE_STORAGE", "http://sample-storage")
	repoSync = func(syncObj common.SyncOps, remoteStorage string) error {
		return errors.New("test-repo-sync-error")
	}
	err := RepoSync()
	assert.EqualError(t, err, "test-repo-sync-error")
}

//func TestGitSync(t *testing.T) {
//	getSyncAuth = func(authObj AuthOps, remoteStorage string) (transport.AuthMethod, error) {
//		auth := &http.BasicAuth{Username: "test-user", Password: "test-password"}
//		return auth, nil
//	}
//	testSyncObj := new(mockSync)
//	if0Dir = "config"
//	testSyncObj.On("GitInit").Return(&git.Repository{}, nil)
//	testSyncObj.On("addRemote").Return(nil)
//	testSyncObj.On("open").Return(&git.Repository{}, nil)
//	testSyncObj.On("pull").Return(&git.Worktree{}, nil)
//	testSyncObj.On("status").Return(git.Status{}, nil)
//	testSyncObj.On("addFile").Return(nil)
//	testSyncObj.On("commit").Return(nil)
//	testSyncObj.On("push").Return(nil)
//	err := gitSync(testSyncObj, "http://sample-storage")
//	assert.Nil(t, err)
//}