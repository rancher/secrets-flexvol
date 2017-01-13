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
		bulkSecrets: &bulkSecret{
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
					Mode:       "0777",
				},
			},
		},
	}

	expectedValues = map[string]string{
		"database_password": "hello",
		"database_username": "hello",
	}

	paramFixture = map[string]interface{}{
		"io.rancher.secrets.request_token": "onetime",
	}
)

type testDecryptor struct {
	key *bulkSecret
}

type testGetter struct {
	bulkSecrets *bulkSecret
}

func (td testDecryptor) Decrypt(text string) ([]byte, error) {
	key, err := loadPrivateKeyFromString(insecureKey)
	if err != nil {
		return []byte{}, err
	}

	return rsaDecrypt(key, text)
}

func (tg testGetter) GetSecrets(params map[string]interface{}) (*bulkSecret, error) {
	return tg.bulkSecrets, nil
}

func TestWriter(t *testing.T) {
	dstDir := "/tmp/testdata"

	sw, err := NewRSASecretFileWriter(testDecryptor{}, paramFixture)
	if err != nil {
		t.Error(err)
		return
	}

	// Setups the temp directory
	if err = os.MkdirAll(dstDir, 0755); err != nil {
		t.Error(err)
	}

	secrets, _ := tGet.GetSecrets(map[string]interface{}{})

	// Calls write method
	if err = sw.Write(secrets, dstDir); err != nil {
		t.Error(err)
		return
	}

	//verifies writes happened
	for _, secret := range secrets.Data {
		secret.setDefaults()

		fi, err := os.Stat(path.Join(dstDir, secret.Name))
		if err != nil {
			t.Error(err)
		}
		mode, _ := strconv.ParseUint(secret.Mode, 0, 32)
		if fi.Mode() != os.FileMode(mode) {
			t.Errorf("Mode not set correctly, expected: %d got %d", os.FileMode(mode), fi.Mode())
		}

		uid, _ := strconv.ParseUint(secret.UID, 0, 32)
		if uint64(fi.Sys().(*syscall.Stat_t).Uid) != uid {
			t.Errorf("UID not set correctly, expected: %d got %d", uid, fi.Sys().(*syscall.Stat_t).Uid)
		}

		gid, _ := strconv.ParseUint(secret.GID, 0, 32)
		if uint64(fi.Sys().(*syscall.Stat_t).Gid) != gid {
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