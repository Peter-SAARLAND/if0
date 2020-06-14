package sync

import (
	"errors"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/ssh"
	"if0/common"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type mockAuth struct {
	mock.Mock
}

func (mAuth *mockAuth) parseSSHKey(sshKey []byte) (ssh.Signer, error) {
	args := mAuth.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(ssh.Signer), args.Error(1)
	}
}

func (mAuth *mockAuth) readPassword() ([]byte, error) {
	args := mAuth.Called()
	return []byte(args.String(0)), args.Error(1)
}

func (mAuth *mockAuth) parseSSHKeyWithPassphrase(sshKey, passphrase []byte) (ssh.Signer, error) {
	args := mAuth.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(ssh.Signer), args.Error(1)
	}
}

func TestGetAuth(t *testing.T) {
	testObj := new(mockAuth)
	testObj.On("readPassword").Return("test-user", nil).Once()
	testObj.On("readPassword").Return("test-password", nil).Once()
	auth, err := getAuth(testObj, "http://sample-storage")
	expectedAuth := &http.BasicAuth{Username: "test-user", Password: "test-password"}
	assert.Equal(t, auth, expectedAuth)
	assert.Nil(t, err)
}

func TestGetAuthUsernameError(t *testing.T) {
	testObj := new(mockAuth)
	testObj.On("readPassword").Return("", errors.New("test-error")).Once()
	auth, err := getAuth(testObj, "http://sample-storage")
	assert.Equal(t, auth, nil)
	assert.EqualError(t, err, "test-error")
}

func TestGetAuthPasswordError(t *testing.T) {
	testObj := new(mockAuth)
	testObj.On("readPassword").Return("test-user", nil).Once()
	testObj.On("readPassword").Return("", errors.New("test-error")).Once()
	auth, err := getAuth(testObj, "http://sample-storage")
	assert.Equal(t, auth, nil)
	assert.EqualError(t, err, "test-error")
}

func TestGetSSHAuthNoFile(t *testing.T) {
	common.RootPath = filepath.Join("config")
	authObj := Auth{}
	auth, err := getAuth(&authObj, "git@gitlab:sample-storage")
	assert.Nil(t, auth)
	assert.NotNil(t, err)
}

func TestGetAuthSSH(t *testing.T) {
	common.RootPath = "testdata"
	sshKey, err := ioutil.ReadFile("testdata/.ssh/id_rsa")
	signer, _ := ssh.ParsePrivateKeyWithPassphrase(sshKey, []byte("pwd"))

	testObj := new(mockAuth)
	testObj.On("readPassword").Return("pwd", nil)
	testObj.On("parseSSHKey").Return(nil, errors.New("ssh: this private key is passphrase protected"))
	testObj.On("parseSSHKeyWithPassphrase").Return(signer, nil)
	auth, err := getAuth(testObj, "git@gitlab:sample-storage")
	//expectedAuth := &gitssh.PublicKeys{User: "git", Signer: signer}
	expectedAuth := &gitssh.PublicKeys{User: "git", Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}}
	assert.ObjectsAreEqual(auth, expectedAuth)
	assert.Nil(t, err)
}

func TestGetAuthSSHPassphraseError(t *testing.T) {
	testObj := new(mockAuth)
	testObj.On("parseSSHKey").Return(nil, errors.New("ssh: this private key is passphrase protected"))
	testObj.On("readPassword").Return("", errors.New("test-error"))
	auth, err := getAuth(testObj, "git@gitlab:sample-storage")
	assert.Nil(t, auth)
	assert.EqualError(t, err, "test-error")
}

func TestGetAuthSSHParseError(t *testing.T) {
	testObj := new(mockAuth)
	testObj.On("parseSSHKey").Return(nil, errors.New("test-parse-error"))
	auth, err := getAuth(testObj, "git@gitlab:sample-storage")
	assert.Nil(t, auth)
	assert.EqualError(t, err, "test-parse-error")
}

func TestGetSSHAuthError(t *testing.T) {
	var testObj Auth
	auth, err := getAuth(&testObj, "sample-storage")
	assert.Nil(t, auth)
	assert.EqualError(t, err, "invalid url")
}

func TestGetAuthSSHParseWithPassphraseError(t *testing.T) {
	testObj := new(mockAuth)
	testObj.On("parseSSHKey").Return(nil, errors.New("ssh: this private key is passphrase protected"))
	testObj.On("readPassword").Return("test-passphrase", nil)
	testObj.On("parseSSHKeyWithPassphrase").Return(nil, errors.New("test-parse-error"))
	auth, err := getAuth(testObj, "git@gitlab:sample-storage")
	assert.Nil(t, auth)
	assert.EqualError(t, err, "test-parse-error")
}

func TestParseGitConfig(t *testing.T) {
	common.RootPath = "testdata"
	user, email := getUserConfig()
	assert.Equal(t, user, "if0")
	assert.Equal(t, email, "if0@if0.com")
}
