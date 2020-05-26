package environments

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"if0/common"
	"if0/common/sync"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func cloneEmptyRepo(remoteStorage string) error {
	syncObj := sync.Sync{}
	dirName := strings.Split(filepath.Base(remoteStorage), ".")[0]
	dirPath := filepath.Join(common.EnvDir, dirName)
	r, err := syncObj.GitInit(dirPath)
	if err != nil {
		return err
	}
	// git remote add <repo>
	err = syncObj.AddRemote(remoteStorage, r)
	if err != nil {
		return err
	}
	return nil
}

// This function checks if the environment directory contains necessary files, if not, creates them.
func envInit(envName string) {
	envPath := filepath.Join(common.EnvDir, envName)
	createFile(filepath.Join(envPath, "zero.env"))
	createFile(filepath.Join(envPath, "dash1.env"))
	f := createFile(filepath.Join(envPath, ".gitlab-ci.yml"))
	if f != nil {
		dataToWrite := "include:\n  - remote: 'https://gitlab.com/peter.saarland/scratch/-/raw/master/ci/templates/shipmate.gitlab-ci.yml'"
		_, _ = f.Write([]byte(dataToWrite))
	}
	sshDir := filepath.Join(envPath, ".ssh")
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		_ = os.Mkdir(sshDir, os.ModeDir)
		err := generateSSHKeyPair(sshDir)
		if err != nil {
			fmt.Println("Error: Generating SSH Key pair - ", err)
		}
	}
}

func createFile(fileName string) *os.File {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
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
