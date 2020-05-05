package sync

import (
	"errors"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/ssh"
	"testing"
)

type mockAuth struct {
	mock.Mock
}

func (mAuth *mockAuth) parseSSHKey(sshKey []byte) (ssh.Signer, error) {
	panic("implement me")
}

func (mAuth *mockAuth) readPassword() ([]byte, error) {
	args := mAuth.Called()
	return []byte(args.String(0)), args.Error(1)
}

func (mAuth *mockAuth) parseSSHKeyWithPassphrase(sshKey, passphrase []byte) (ssh.Signer, error) {
	args := mAuth.Called()
	return nil, args.Error(1)
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

//func TestGetSSHAuthNoFile(t *testing.T) {
//	RootPath = filepath.Join("config")
//	authObj := Auth{}
//	Auth, err := getAuth(&authObj, "git@gitlab:sample-storage")
//	assert.Nil(t, Auth)
//	assert.NotNil(t, err)
//}
//
//func TestGetAuthSSH(t *testing.T) {
//	RootPath = filepath.Join("config")
//	testObj := new(mockAuth)
//	_ = os.RemoveAll(".ssh")
//	generateTestPrivKey()
//	fmt.Println(filepath.Join(RootPath, ".ssh", "id_rsa"))
//	testObj.On("readPassword").Return("pwd", nil)
//	Auth, err := getAuth(testObj, "git@gitlab:sample-storage")
//	expectedAuth := &gitssh.PublicKeys{User: "git", Signer: nil}
//	assert.Equal(t, Auth, expectedAuth)
//	assert.Nil(t, err)
//	_ = os.RemoveAll(".ssh")
//}
//
//func TestGetAuthSSHPassphraseError(t *testing.T) {
//	testObj := new(mockAuth)
//	testObj.On("readPassword").Return("", errors.New("test-error"))
//	Auth, err := getAuth(testObj, "git@gitlab:sample-storage")
//	assert.Nil(t, Auth)
//	assert.EqualError(t, err, "test-error")
//}
//
//func TestGetAuthSSHParseError(t *testing.T) {
//	testObj := new(mockAuth)
//	testObj.On("readPassword").Return("test-passphrase", nil)
//	testObj.On("parseSSHKeyWithPassphrase").Return(nil, errors.New("test-parse-error"))
//	Auth, err := getAuth(testObj, "git@gitlab:sample-storage")
//	assert.Nil(t, Auth)
//	assert.EqualError(t, err, "test-parse-error")
//}
//
//func generateTestPrivKey() {
//	RootPath = filepath.Join("config")
//	_ = os.Mkdir(".ssh", 0777)
//	f, err := os.OpenFile(filepath.Join(".ssh", "id_rsa"), os.O_CREATE|os.O_RDWR, 0777)
//	fmt.Println(err)
//	key, _ := rsa.GenerateKey(rand.Reader, 2048)
//	// Convert it to pem
//	block := &pem.Block{
//		Type:  "RSA PRIVATE KEY",
//		Bytes: x509.MarshalPKCS1PrivateKey(key),
//	}
//	// Encrypt the pem
//	block, _ = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte("pwd"), x509.PEMCipherAES256)
//	pembytes := pem.EncodeToMemory(block)
//	_ = ioutil.WriteFile(f.Name(), pembytes, 0666)
//}