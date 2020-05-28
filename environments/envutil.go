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
	createFile(filepath.Join(envPath, "dash1.env"))
	f := createFile(filepath.Join(envPath, ".gitlab-ci.yml"))
	if f != nil {
		shipmateUrl := config.GetEnvVariable("SHIPMATE_WORKFLOW_URL")
		dataToWrite := fmt.Sprintf("include:\n  - remote: '%s'", shipmateUrl)
		_, _ = f.Write([]byte(dataToWrite))
	}
	sshDir := filepath.Join(envPath, ".ssh")
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		_ = os.Mkdir(sshDir, os.ModeDir)
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
	fmt.Println("Syncing environment init file changes")
	w, _ := syncObj.GetWorktree(r)
	status, _ := syncObj.Status(w)
	if len(status) > 0 {
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

	err = writeKeyToFile(privateKeyBytes, privateKeyPath)
	if err != nil {
		return err
	}

	err = writeKeyToFile([]byte(publicKeyBytes), publicKeyPath)
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

func writeKeyToFile(keyBytes []byte, file string) error {
	err := ioutil.WriteFile(file, keyBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}
