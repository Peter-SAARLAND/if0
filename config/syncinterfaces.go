package config

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type syncOps interface {
	gitInit(localRepoPath string) (*git.Repository, error)
	addRemote(remoteStorage string, r *git.Repository) error
	open(if0Dir string) (*git.Repository, error)
	pull(remoteStorage string, r *git.Repository, pullOptions *git.PullOptions) (*git.Worktree, error)
	status(w *git.Worktree) (git.Status, error)
	addFile(w *git.Worktree, file string) error
	commit(w *git.Worktree) error
	push(auth transport.AuthMethod, r *git.Repository) error
}

type sync struct {
}

func (s *sync) gitInit(localRepoPath string) (*git.Repository, error) {
	log.Println("Creating a git repository at ", localRepoPath)
	r, err := git.PlainInit(localRepoPath, false)
	if err != nil {
		log.Errorln("Error while creating a git repository: ", err)
		return nil, err
	}
	return r, nil
}

func (s *sync) addRemote(remoteStorage string, r *git.Repository) error {
	log.Println("Adding remote 'origin' for the repository at ", remoteStorage)
	remoteConfig := &config.RemoteConfig{Name: "origin", URLs: []string{remoteStorage}}
	_, err := r.CreateRemote(remoteConfig)
	if err != nil {
		log.Errorln("Error while adding remote: ", err)
		return err
	}
	return nil
}

func (s *sync) open(if0Dir string) (*git.Repository, error) {
	log.Println("Git repository already present.")
	return git.PlainOpen(if0Dir)
}

func (s *sync) pull(remoteStorage string, r *git.Repository, pullOptions *git.PullOptions) (*git.Worktree, error) {
	log.Println("Pulling in changes from ", remoteStorage)
	w, err := r.Worktree()
	if err != nil {
		log.Errorln("Worktree error: ", err)
		return nil, err
	}
	err = w.Pull(pullOptions)
	if err != nil {
		log.Errorln("Pull status: ", err)
		return w, err
	}
	return w, nil
}

func (s *sync) status(w *git.Worktree) (git.Status, error) {
	return w.Status()
}

func (s *sync) addFile(w *git.Worktree, file string) error {
	_, err := w.Add(file)
	return err
}

func (s *sync) commit(w *git.Worktree) error {
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

func (s *sync) push(auth transport.AuthMethod, r *git.Repository) error {
	log.Println("Pushing local changes")
	pushOptions := &git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	}
	return r.Push(pushOptions)
}

