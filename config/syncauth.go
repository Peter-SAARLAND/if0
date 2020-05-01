package config

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	getSyncAuth = getAuth
)

type Auth interface {
	readPassword() ([]byte, error)
	parseSSHKeyWithPassphrase(sshKey, passphrase []byte) (ssh.Signer, error)
	parseSSHKey(sshKey []byte) (ssh.Signer, error)
}

type auth struct {
}

func getAuth(authObj Auth, remoteStorage string) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	var err error
	if strings.Contains(remoteStorage, "http") {
		auth, err = getHttpAuth(authObj)
		if err != nil {
			log.Errorln("Error while fetching credentials: ", err)
			return nil, err
		}
	} else if strings.Contains(remoteStorage, "git@") {
		auth, err = getSSHAuth(authObj)
		if err != nil {
			log.Errorln("Error while fetching credentials: ", err)
			return nil, err
		}
	}
	return auth, nil
}

func getHttpAuth(authObj Auth) (transport.AuthMethod, error) {
	fmt.Println("Enter Username: ")
	userName, err := authObj.readPassword()
	if err != nil {
		fmt.Printf("Failed to read username: %v", err)
		return nil, err
	}
	fmt.Println("Enter Password: ")
	bytePassword, err := authObj.readPassword()
	if err != nil {
		fmt.Printf("Failed to read password: %v", err)
		return nil, err
	}
	auth := &http.BasicAuth{Username: string(userName), Password: string(bytePassword)}
	return auth, nil
}

func getSSHAuth(authObj Auth) (*gitssh.PublicKeys, error) {
	sshKeyPath := filepath.Join(rootPath, ".ssh", "id_rsa")
	sshKey, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		fmt.Println("Error while reading SSH key: ", err)
		return nil, err
	}
	signer, err := authObj.parseSSHKey(sshKey)
	if err != nil {
		if err.Error() == "ssh: this private key is passphrase protected" {
			fmt.Println("Passphrase required. Enter Passphrase")
			passphrase, err := authObj.readPassword()
			if err != nil {
				log.Println("Error while reading passphrase: ", err)
				return nil, err
			}
			signer, err = authObj.parseSSHKeyWithPassphrase(sshKey, passphrase)
			if err != nil {
				fmt.Println("Error while parsing SSH key: ", err)
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	auth := &gitssh.PublicKeys{User: "git", Signer: signer}
	return auth, nil
}

func (p *auth) parseSSHKeyWithPassphrase(sshKey, passphrase []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKeyWithPassphrase(sshKey, passphrase)

}

func (p *auth) parseSSHKey(sshKey []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKey(sshKey)
}

func (p *auth) readPassword() ([]byte, error) {
	secret, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	return secret, nil
}