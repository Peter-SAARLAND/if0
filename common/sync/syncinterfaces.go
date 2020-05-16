package sync

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"if0/common"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type SyncOps interface {
	GitInit(localRepoPath string) (*git.Repository, error)
	AddRemote(remoteStorage string, r *git.Repository) error
	Open(if0Dir string) (*git.Repository, error)
	Pull(remoteStorage string, r *git.Repository, pullOptions *git.PullOptions) (*git.Worktree, error)
	Status(w *git.Worktree) (git.Status, error)
	AddFile(w *git.Worktree, file string) error
	Commit(w *git.Worktree) error
	Push(auth transport.AuthMethod, r *git.Repository) error
	Clone(repoUrl string, auth transport.AuthMethod) (*git.Repository, error)
	GetWorktree(r *git.Repository) (*git.Worktree, error)
}

type Sync struct {
}

func (s *Sync) GitInit(localRepoPath string) (*git.Repository, error) {
	fmt.Println("Creating a git repository at ", localRepoPath)
	r, err := git.PlainInit(localRepoPath, false)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Sync) AddRemote(remoteStorage string, r *git.Repository) error {
	fmt.Println("Adding remote 'origin' for the repository at ", remoteStorage)
	remoteConfig := &config.RemoteConfig{Name: "origin", URLs: []string{remoteStorage}}
	_, err := r.CreateRemote(remoteConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sync) Open(if0Dir string) (*git.Repository, error) {
	return git.PlainOpen(if0Dir)
}

func (s *Sync) Pull(remoteStorage string, r *git.Repository, pullOptions *git.PullOptions) (*git.Worktree, error) {
	fmt.Println("Pulling in changes from ", remoteStorage)
	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}
	err = w.Pull(pullOptions)
	if err != nil {
		return w, err
	}
	return w, nil
}

func (s *Sync) Status(w *git.Worktree) (git.Status, error) {
	return w.Status()
}

func (s *Sync) AddFile(w *git.Worktree, file string) error {
	_, err := w.Add(file)
	return err
}

func (s *Sync) Commit(w *git.Worktree) error {
	name, email := getUserConfig()
	commitMsg := "feat: updating config files - " + time.Now().Format("02012006_150405")
	commitOptions := &git.CommitOptions{
		All: false,
		Author: &object.Signature{
			When: time.Now(),
			Name: name,
			Email: email,
		},
	}
	_, err := w.Commit(commitMsg, commitOptions)
	return err
}

func (s *Sync) Push(auth transport.AuthMethod, r *git.Repository) error {
	fmt.Println("Pushing local changes")
	pushOptions := &git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	}
	return r.Push(pushOptions)
}

func (s *Sync) Clone(repoUrl string, auth transport.AuthMethod) (*git.Repository, error) {
	fmt.Printf("Cloning the git repository %s at %s\n", repoUrl, common.EnvDir)
	cloneOptions := &git.CloneOptions{
		URL:      repoUrl,
		Auth:     auth,
		Progress: os.Stdout,
	}

	localRepoPath := filepath.Join(common.EnvDir, strings.Split(path.Base(repoUrl), ".")[0])
	r, err := git.PlainClone(localRepoPath, false, cloneOptions)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Sync) GetWorktree(r *git.Repository) (*git.Worktree, error) {
	return r.Worktree()
}