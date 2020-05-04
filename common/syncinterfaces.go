package common

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	log "github.com/sirupsen/logrus"
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
}

type Sync struct {
}

func (s *Sync) GitInit(localRepoPath string) (*git.Repository, error) {
	log.Println("Creating a git repository at ", localRepoPath)
	r, err := git.PlainInit(localRepoPath, false)
	if err != nil {
		log.Errorln("Error while creating a git repository: ", err)
		return nil, err
	}
	return r, nil
}

func (s *Sync) AddRemote(remoteStorage string, r *git.Repository) error {
	log.Println("Adding remote 'origin' for the repository at ", remoteStorage)
	remoteConfig := &config.RemoteConfig{Name: "origin", URLs: []string{remoteStorage}}
	_, err := r.CreateRemote(remoteConfig)
	if err != nil {
		log.Errorln("Error while adding remote: ", err)
		return err
	}
	return nil
}

func (s *Sync) Open(if0Dir string) (*git.Repository, error) {
	log.Println("Git repository already present.")
	return git.PlainOpen(if0Dir)
}

func (s *Sync) Pull(remoteStorage string, r *git.Repository, pullOptions *git.PullOptions) (*git.Worktree, error) {
	log.Println("Pulling in changes from ", remoteStorage)
	w, err := r.Worktree()
	if err != nil {
		log.Errorln("Worktree error: ", err)
		return nil, err
	}
	err = w.Pull(pullOptions)
	if err != nil {
		log.Errorln("Pull Status: ", err)
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
	commitMsg := "feat: updating config files - " + time.Now().Format("02012006_150405")
	commitOptions := &git.CommitOptions{
		All: false,
		Author: &object.Signature{
			When: time.Now(),
		},
	}
	_, err := w.Commit(commitMsg, commitOptions)
	return err
}

func (s *Sync) Push(auth transport.AuthMethod, r *git.Repository) error {
	log.Println("Pushing local changes")
	pushOptions := &git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	}
	return r.Push(pushOptions)
}

func (s *Sync) Clone(repoUrl string, auth transport.AuthMethod) (*git.Repository, error) {
	log.Printf("Cloning the git repository %s at %s", repoUrl, EnvDir)
	cloneOptions := &git.CloneOptions{
		URL:      repoUrl,
		Auth:     auth,
		Progress: os.Stdout,
	}

	localRepoPath := filepath.Join(EnvDir, strings.Split(path.Base(repoUrl), ".")[0])
	r, err := git.PlainClone(localRepoPath, false, cloneOptions)
	if err != nil {
		return nil, err
	}
	return r, nil
}
