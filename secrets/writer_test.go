package secrets

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"syscall"
	"testing"
)

var (
	tGet = &testGetter{
		Data: []secret{
			secret{
				Name:       "database_password",
				RewrapText: "Yyaba6uZYkPHLqzzh4n6SB76tU32ugonB8uxdViUhxKpk/tThhPdQQvj4pe1k3advNOMUyIuykbnJ9EUVY4M4KRdilt6KlCQTEPrzTGw9ZxoFdBWlW2Kj3+1BZt/iy36krzryyLS+bNIDE8IRNoafaPmcto1ywQHfBjXiIjoJfYIuXpbQPOLU1ulElMv7ArwG2JbIvYcpIMysoJqaJ7YAauHveMPmAbRB/oGgS/pxIoP9vv1PMPIoP6c6h4raWXZ6uRkMJ7ND6cEq3pXLVlapYgZnOV9lbMBxQGlzApVlDo4BnMsNz/NNiaKYQs5CjO12KySuDjLkRamERL1FaKQhA==",
			},
			secret{
				Name:       "database_username",
				RewrapText: "Yyaba6uZYkPHLqzzh4n6SB76tU32ugonB8uxdViUhxKpk/tThhPdQQvj4pe1k3advNOMUyIuykbnJ9EUVY4M4KRdilt6KlCQTEPrzTGw9ZxoFdBWlW2Kj3+1BZt/iy36krzryyLS+bNIDE8IRNoafaPmcto1ywQHfBjXiIjoJfYIuXpbQPOLU1ulElMv7ArwG2JbIvYcpIMysoJqaJ7YAauHveMPmAbRB/oGgS/pxIoP9vv1PMPIoP6c6h4raWXZ6uRkMJ7ND6cEq3pXLVlapYgZnOV9lbMBxQGlzApVlDo4BnMsNz/NNiaKYQs5CjO12KySuDjLkRamERL1FaKQhA==",
				UID:        "1000",
				GID:        "1000",
				Mode:       "777",
			},
		},
	}

	expectedValues = map[string]string{
		"database_password": "hello",
		"database_username": "hello",
	}

	paramFixture = &options{
		Token: &secretToken{
			Value: []byte("onetime"),
		},
	}
)

type testDecryptor struct{}

type testGetter struct {
	Data []secret
}

func (td testDecryptor) Decrypt(text string) ([]byte, error) {
	key, err := loadPrivateKeyFromString(insecureKey)
	if err != nil {
		return []byte{}, err
	}

	return rsaDecrypt(key, text)
}

func (tg testGetter) GetSecrets(params *options) ([]secret, error) {
	return tg.Data, nil
}

func TestWriter(t *testing.T) {
	dstDir := "/tmp/testdata"

	sw, err := NewRSASecretFileWriter(testDecryptor{})
	if err != nil {
		t.Error(err)
		return
	}

	// Setups the temp directory
	if err = os.MkdirAll(dstDir, 0755); err != nil {
		t.Error(err)
	}

	secrets, _ := tGet.GetSecrets(paramFixture)

	// Calls write method
	if err = sw.Write(secrets, dstDir); err != nil {
		t.Error(err)
		return
	}

	//verifies writes happened
	for _, secret := range secrets {
		secret.setDefaults()

		fi, err := os.Stat(path.Join(dstDir, secret.Name))
		if err != nil {
			t.Error(err)
		}
		mode, _ := strconv.ParseUint(secret.Mode, 8, 32)
		if fi.Mode() != os.FileMode(mode) {
			t.Errorf("Mode not set correctly, expected: %d got %d", os.FileMode(mode), fi.Mode())
		}

		uid, _ := strconv.Atoi(secret.UID)
		if int(fi.Sys().(*syscall.Stat_t).Uid) != uid {
			t.Errorf("UID not set correctly, expected: %d got %d", uid, fi.Sys().(*syscall.Stat_t).Uid)
		}

		gid, _ := strconv.Atoi(secret.GID)
		if int(fi.Sys().(*syscall.Stat_t).Gid) != gid {
			t.Errorf("Mode not set correctly, expected: %d got %d", uid, fi.Sys().(*syscall.Stat_t).Gid)
		}

		content, err := ioutil.ReadFile(path.Join(dstDir, secret.Name))
		if err != nil {
			t.Error(err)
		}

		if expectedValues[secret.Name] != string(content) {
			t.Errorf("Contents of file %s not equal. %s != %s", fi.Name(), expectedValues[secret.Name], string(content))
		}

	}
	return
}
