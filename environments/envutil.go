package environments

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"if0/common"
	"if0/common/sync"
	"if0/config"
	gitlabclient "if0/environments/git"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	pushEnvInitChanges = pushInitChanges
)

func cloneEnv(repoUrl, envDir string) (*git.Repository, error) {
	// get authorization
	authObj := sync.Auth{}
	auth, err := getAuth(&authObj, repoUrl)
	if err != nil {
		fmt.Println("Authentication error - ", err)
		return nil, err
	}

	r, err := clone(repoUrl, envDir, auth)
	if err != nil {
		if err.Error() == "remote repository is empty" {
			r, err = cloneEmptyRepo(repoUrl, envDir)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return r, nil
}

func cloneEmptyRepo(remoteStorage, envDir string) (*git.Repository, error) {
	syncObj := sync.Sync{}
	//dirName := strings.Split(filepath.Base(remoteStorage), ".")[0]
	//dirPath := filepath.Join(common.EnvDir, dirName)
	r, err := syncObj.GitInit(envDir)
	if err != nil {
		return nil, err
	}
	// git remote add <repo>
	err = syncObj.AddRemote(remoteStorage, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// This function checks if the environment directory contains necessary files, if not, creates them.
func envInit(envPath string) error {
	createZeroFile(envPath)
	createCIFile(envPath)
	sshDir := filepath.Join(envPath, ".ssh")
	files, direrr := ioutil.ReadDir(sshDir)
	// .ssh dir not present or present but no keys
	if _, err := os.Stat(sshDir); os.IsNotExist(err) || (direrr == nil && len(files) < 2) {
		fmt.Printf("Creating dir %s\n", sshDir)
		_ = os.Mkdir(sshDir, 0700)
		err := generateSSHKeyPair(sshDir)
		if err != nil {
			fmt.Println("Error: Generating SSH Key pair - ", err)
			return err
		}
	}
	return nil
}

func createZeroFile(envPath string) {
	f := createFile(filepath.Join(envPath, "zero.env"))
	defer f.Close()
	pwd := generateRandSeq()
	hash, err := generateHashCmd(pwd)
	if runtime.GOOS == "windows" || hash == "" || err != nil {
		hash, err = generateHashDocker(pwd)
		if err != nil {
			fmt.Println("Error: Could not create htpasswd hash -", err)
			return
		}
	}
	_, _ = f.WriteString("ZERO_ADMIN_USER=admin\n")
	_, _ = f.WriteString("ZERO_ADMIN_PASSWORD="+pwd+"\n")
	_, _ = f.WriteString("ZERO_ADMIN_PASSWORD_HASH="+hash+"\n")
}

func createCIFile(envPath string) {
	f := createFile(filepath.Join(envPath, ".gitlab-ci.yml"))
	defer f.Close()
	if f != nil {
		shipmateUrl := getShipmateUrl()
		dataToWrite := fmt.Sprintf("include:\n  - remote: '%s'", shipmateUrl)
		_, _ = f.Write([]byte(dataToWrite))
	}
}

func pushInitChanges(r *git.Repository, auth transport.AuthMethod) error {
	w, _ := syncObj.GetWorktree(r)
	status, _ := syncObj.Status(w)
	if len(status) > 0 {
		fmt.Println("Syncing environment init file changes")
		for file, _ := range status {
			_ = syncObj.AddFile(w, file)
		}
		// git commit
		err := syncObj.Commit(w)
		if err != nil {
			fmt.Println("Error: Committing changes - ", err)
			return err
		}
		// git push
		err = syncObj.Push(auth, r)
		if err != nil {
			fmt.Println("Error: Pushing changes - ", err)
			return err
		}
	}
	return nil
}

func createFile(fileName string) *os.File {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("Creating file", fileName)
		f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
		return f
	}
	return nil
}

func getShipmateUrl() string {
	// get SHIPMATE_WORKFLOW_URL from if0.env
	config.ReadConfigFile(common.If0Default)
	shipmateUrl := config.GetEnvVariable("SHIPMATE_WORKFLOW_URL")
	// if not found, add it to if0.env and return the value
	if shipmateUrl == "" {
		f, _ := os.OpenFile(common.If0Default, os.O_APPEND, 0644)
		defer f.Close()
		_, _ = f.WriteString("SHIPMATE_WORKFLOW_URL=https://gitlab.com/peter.saarland/shipmate/-/raw/master/shipmate.gitlab-ci.yml\n")
	}
	config.ReadConfigFile(common.If0Default)
	return config.GetEnvVariable("SHIPMATE_WORKFLOW_URL")
}

func createLocalEnv(repoName string, repoUrl string) error {
	envDir := createNestedDirPath(repoName, repoUrl)
	// if a remote repository (empty) url is provided, sync the changes
	if repoUrl != "" {
		_, _ = cloneEnv(repoUrl, envDir)
		addLocalEnv(envDir)
		err := syncLocalEnvChanges(repoUrl, envDir)
		if err != nil {
			return err
		}
	} else {
		addLocalEnv(envDir)
		fmt.Println("No remote repository url was found for sync. "+
			"The environment has been created locally at ", envDir)
		fmt.Println("To sync the local changes, run `if0 add repo-name repo-url`")
	}
	return nil
}

func createGLProject(repoName, glToken string) error {
	// creating a private project in gitlab
	sshRepoUrl, err := gitlabclient.CreateProject(repoName, glToken)
	if err != nil {
		return err
	}
	// adding the environment locally
	envDir := createNestedDirPath(repoName, sshRepoUrl)
	addLocalEnv(envDir)
	// syncing local changes with the private project
	err = syncLocalEnvChanges(sshRepoUrl, envDir)
	if err != nil {
		return err
	}
	return nil
}

func syncLocalEnvChanges(repoUrl string, envDir string) error {
	authObj := sync.Auth{}
	auth, err := getAuth(&authObj, repoUrl)
	if err != nil {
		fmt.Println("Authentication error - ", err)
		return err
	}
	r, _ := config.GetRepository(&syncObj, repoUrl, envDir)
	err = pushEnvInitChanges(r, auth)
	if err != nil {
		fmt.Println("Error: Pushing env init changes -", err)
		return err
	}
	return nil
}

func addLocalEnv(envDir string) {
	// check if the repo exists already.
	// if it does not exist, create a new one locally and sync
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		err = os.MkdirAll(envDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error: Creating nested directories - ", err)
		}
		_ = envInit(envDir)
	}
}

func createNestedDirPath(repoName, repoUrl string) string {
	var dirPath string
	if repoUrl != "" {
		dirPathElem := strings.FieldsFunc(repoUrl, func(r rune) bool {
			return r == ':' || r == '/' || r == '@'
		})
		dirPathElem[len(dirPathElem)-1] = strings.Split(dirPathElem[len(dirPathElem)-1], ".")[0]
		dirPath = filepath.Join(common.EnvDir, strings.Join(dirPathElem[1:], string(os.PathSeparator)))
	} else {
		dirPath = filepath.Join(common.EnvDir, repoName)
	}
	return dirPath
}