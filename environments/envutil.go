package environments

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"if0/common"
	"if0/common/sync"
	"if0/config"
	gitlabclient "if0/environments/git"
	"net/url"
	"os"
	"path/filepath"
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
			"The local copy of the environment can be found at ", envDir)
		fmt.Println("To sync the local changes, run `if0 add repo-name repo-url`")
	}
	return nil
}

func createGLProject(repoName, glToken string) error {
	config.ReadConfigFile(common.If0Default)
	if0RegUrl := config.GetEnvVariable("IF0_REGISTRY_URL")
	if0RegGroup := config.GetEnvVariable("IF0_REGISTRY_GROUP")
	httpRepoUrl := if0RegUrl+"/"+if0RegGroup+"/"+repoName
	// adding the environment locally
	// TODO: we need a check here to check if the project exists already on gitlab
	envDir := createNestedDirPath(repoName, httpRepoUrl)
	addLocalEnv(envDir)
	// creating a private project in gitlab
	sshRepoUrl, _, err := gitlabclient.CreateProject(repoName, glToken)
	if err != nil {
		return err
	}
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
	} else {
		fmt.Println("environment already present")
		return
	}
	_ = envInit(envDir)
}

func createNestedDirPath(repoName, repoUrl string) string {
	var dirPath string
	if repoUrl != "" {
		repoUrl = removePortNumber(repoUrl)
		dirPathElem := strings.FieldsFunc(repoUrl, func(r rune) bool {
			return r == ':' || r == '/' || r == '@'
		})
		dirPathElem[len(dirPathElem)-1] = strings.Split(dirPathElem[len(dirPathElem)-1], ".")[0]
		dirPath = filepath.Join(common.EnvDir, strings.Join(dirPathElem, string(os.PathSeparator)))
	} else {
		dirPath = filepath.Join(common.EnvDir, repoName)
	}
	return dirPath
}

func removePortNumber(repoUrl string) string {
	if strings.Contains(repoUrl, "ssh:") {
		u, _ := url.Parse(repoUrl)
		repoUrl = strings.Replace(repoUrl, u.Port(), "", 1)
		i := strings.Index(repoUrl, "git@")
		repoUrl = repoUrl[i+4:]
	} else if strings.Contains(repoUrl, "git@") {
		i := strings.Index(repoUrl, "git@")
		repoUrl = repoUrl[i+4:]
	} else if strings.Contains(repoUrl, "http") {
		u, _ := url.Parse(repoUrl)
		repoUrl = strings.Replace(repoUrl, u.Port(), "", 1)
		repoUrl = strings.Replace(repoUrl, u.Scheme, "", 1)
	}
	return repoUrl
}

// helper functions for `if0 list`
func visit(p string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() && (info.Name() == ".git" || info.Name() == ".ssh") {
		return filepath.SkipDir
	}
	hasZeroEnv := checkForZeroEnv(p)
	if info.IsDir() && hasZeroEnv {
		fmt.Printf("Local repository: %s. ", p)
		repoUrl := getRepoUrl(p)
		if repoUrl == "" {
			repoUrl = "remote repository does not exist"
		}
		fmt.Println("Repository URL: ", repoUrl)
	}
	return nil
}

func checkForZeroEnv(dir string) bool {
	zeroPath := filepath.Join(dir, "zero.env")
	if _, err := os.Stat(zeroPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func getRepoUrl(envDir string) string {
	// open the existing repo at ~/.if0
	r, err := syncObj.Open(envDir)
	if err != nil {
		return ""
	}
	remotes, _ := r.Remote("origin")
	return remotes.Config().URLs[0]
}
