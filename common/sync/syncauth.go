package sync

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"if0/common"
	"io/ioutil"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	GetSyncAuth = getAuth
)

type AuthOps interface {
	readPassword() ([]byte, error)
	parseSSHKeyWithPassphrase(sshKey, passphrase []byte) (ssh.Signer, error)
	parseSSHKey(sshKey []byte) (ssh.Signer, error)
}

type Auth struct {
}

func getAuth(authObj AuthOps, remoteStorage string) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	var err error
	if strings.Contains(remoteStorage, "http") {
		auth, err = getHttpAuth(authObj)
		if err != nil {
			fmt.Println("Error: HTTP Authorization - ", err)
			return nil, err
		}
	} else if strings.Contains(remoteStorage, "git@") {
		auth, err = getSSHAuth(authObj)
		if err != nil {
			fmt.Println("Error: SSH Authorization - ", err)
			return nil, err
		}
	} else {
		return nil, errors.New("invalid url")
	}
	return auth, nil
}

func getHttpAuth(authObj AuthOps) (transport.AuthMethod, error) {
	fmt.Println("Enter Username: ")
	userName, err := authObj.readPassword()
	if err != nil {
		fmt.Println("Error: Reading username - ", err)
		return nil, err
	}
	fmt.Println("Enter Password: ")
	bytePassword, err := authObj.readPassword()
	if err != nil {
		fmt.Println("Error: Reading password - ", err)
		return nil, err
	}
	auth := &http.BasicAuth{Username: string(userName), Password: string(bytePassword)}
	return auth, nil
}

func getSSHAuth(authObj AuthOps) (*gitssh.PublicKeys, error) {
	sshKeyPath := filepath.Join(common.RootPath, ".ssh", "id_rsa")
	sshKey, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		fmt.Println("Error: Reading SSH key - ", err)
		return nil, err
	}
	signer, err := authObj.parseSSHKey(sshKey)
	if err != nil {
		if err.Error() == "ssh: this private key is passphrase protected" {
			fmt.Println("Passphrase required. Enter Passphrase")
			passphrase, err := authObj.readPassword()
			if err != nil {
				fmt.Println("Error: Reading passphrase - ", err)
				return nil, err
			}
			signer, err = authObj.parseSSHKeyWithPassphrase(sshKey, passphrase)
			if err != nil {
				fmt.Println("Error: Parsing SSH key - ", err)
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	auth := &gitssh.PublicKeys{User: "git", Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}}
	return auth, nil
}

func (p *Auth) parseSSHKeyWithPassphrase(sshKey, passphrase []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKeyWithPassphrase(sshKey, passphrase)

}

func (p *Auth) parseSSHKey(sshKey []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKey(sshKey)
}

func (p *Auth) readPassword() ([]byte, error) {
	secret, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func getUserConfig() (string, string) {
	var name, email string
	cfg := config.NewConfig()
	gitConfig := filepath.Join(common.RootPath, ".gitconfig")
	b, err := ioutil.ReadFile(gitConfig)
	if err != nil {
		fmt.Println("Error: Reading .gitconfig - ", err)
		return name, email
	}
	err = cfg.Unmarshal(b)
	if err != nil {
		fmt.Println("Error: Unmarshalling .gitconfig - ", err)
		return name, email
	}
	for _, ss := range cfg.Raw.Sections {
		email = ss.Options.Get("email")
		name = ss.Options.Get("name")
		return name, email
	}
	return name, email
}