package environments

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"golang.org/x/crypto/ssh"
	"if0/common"
	"if0/common/sync"
	"if0/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	pushEnvInitChanges = pushInitChanges
)

func cloneEmptyRepo(remoteStorage string) (*git.Repository, error) {
	syncObj := sync.Sync{}
	dirName := strings.Split(filepath.Base(remoteStorage), ".")[0]
	dirPath := filepath.Join(common.EnvDir, dirName)
	r, err := syncObj.GitInit(dirPath)
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
func envInit(r *git.Repository, auth transport.AuthMethod, envName string) error {
	envPath := filepath.Join(common.EnvDir, envName)
	createFile(filepath.Join(envPath, "zero.env"))
	f := createFile(filepath.Join(envPath, ".gitlab-ci.yml"))
	defer f.Close()
	if f != nil {
		shipmateUrl := getShipmateUrl()
		dataToWrite := fmt.Sprintf("include:\n  - remote: '%s'", shipmateUrl)
		_, _ = f.Write([]byte(dataToWrite))
	}
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
	// Pushing the newly added changes to the remote repository
	err := pushEnvInitChanges(r, auth)
	return err
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

func generateSSHKeyPair(sshDir string) error {
	privateKeyPath := filepath.Join(sshDir, "id_rsa")
	publicKeyPath := filepath.Join(sshDir, "id_rsa.pub")

	privateKey, err := generatePrivateKey()
	if err != nil {
		return err
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	err = writeKeyToFile(privateKeyBytes, privateKeyPath, 0600)
	if err != nil {
		return err
	}

	err = writeKeyToFile([]byte(publicKeyBytes), publicKeyPath, 0644)
	if err != nil {
		return err
	}
	return nil
}

func generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func generatePublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	return pubKeyBytes, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privateDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateDER,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	return privatePEM
}

func writeKeyToFile(keyBytes []byte, file string, perm os.FileMode) error {
	fmt.Printf("Creating ssh key %s\n", file)
	err := ioutil.WriteFile(file, keyBytes, perm)
	if err != nil {
		return err
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
