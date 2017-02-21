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
				RewrapText: "eyJlbmNyeXB0aW9uQWxnb3JpdGhtIjoiYWVzMjU2LWdjbTk2IiwiZW5jcnlwdGVkVGV4dCI6IntcIk5vbmNlXCI6XCJMV3QwUGNYN0FjUDQ5ZU1FXCIsXCJBbGdvcml0aG1cIjpcImFlczI1Ni1nY21cIixcIkNpcGhlclRleHRcIjpcImdkU2RLbUMwM0dVK0VWVTlIc1ZzWTFuUTl5NVM5NEJDXCJ9IiwiZW5jcnlwdGVkS2V5Ijp7ImVuY3J5cHRpb25BbGdvcml0aG0iOiJQS0NTMV9PQUVQIiwiZW5jcnlwdGVkVGV4dCI6ImQ5cGcvWEFvY1daRk9zbUF4QytNQzYwT2YxaGdGTnU3UFExRFY5NnJUZzZQRkVQa0x4TGlOK2tyQlFGSENKemNubk1CY3FQUzVGMjdIdWV6dWMzZUFhZVVtUDlmYjJLUlB4b0pRT0tzRnlPbnBDdTd0OFRNRWszWW5MZUZFVTJhWWFmVEhXdnhjYVROM2dBb2xaa2xTOXdRczQ2MzVYcHRxSHBDQkhaN1NJSXJVRWd5Z2t3bEdvbFBJVmxGQUpZQ3dZbGRlTEJsdk9BcHN3WEFmZ05kenVZMThaMEhwSTFhOWRIdko0MGFRRFZ1R21ZY3hSYXMrRnQ1MmRJQ3VvSU5DZWVOZEdwUzd4RkxndWVIaG0xY3JMSjBIRUhCVFpwRno0NWJGSG9vc2kweXoxL0RrNm9RUGhLZStkdmlOUVZocnp5VGtsTHNkZ1dGWW1NNVpVQTMvQT09IiwiaGFzaEFsZ29yaXRobSI6InNoYTI1NiJ9LCJzaWduYXR1cmUiOiI3N3lBWXpCdzR3UmVsN3A0T2lXcWxFZnNXeWFpQVA3VlJpVkhlNnlGWGN5WlArR1Q4UVg1cXdRdGk3TzAifQ==",
			},
			secret{
				Name:       "database_username",
				RewrapText: "eyJlbmNyeXB0aW9uQWxnb3JpdGhtIjoiYWVzMjU2LWdjbTk2IiwiZW5jcnlwdGVkVGV4dCI6IntcIk5vbmNlXCI6XCJMV3QwUGNYN0FjUDQ5ZU1FXCIsXCJBbGdvcml0aG1cIjpcImFlczI1Ni1nY21cIixcIkNpcGhlclRleHRcIjpcImdkU2RLbUMwM0dVK0VWVTlIc1ZzWTFuUTl5NVM5NEJDXCJ9IiwiZW5jcnlwdGVkS2V5Ijp7ImVuY3J5cHRpb25BbGdvcml0aG0iOiJQS0NTMV9PQUVQIiwiZW5jcnlwdGVkVGV4dCI6ImQ5cGcvWEFvY1daRk9zbUF4QytNQzYwT2YxaGdGTnU3UFExRFY5NnJUZzZQRkVQa0x4TGlOK2tyQlFGSENKemNubk1CY3FQUzVGMjdIdWV6dWMzZUFhZVVtUDlmYjJLUlB4b0pRT0tzRnlPbnBDdTd0OFRNRWszWW5MZUZFVTJhWWFmVEhXdnhjYVROM2dBb2xaa2xTOXdRczQ2MzVYcHRxSHBDQkhaN1NJSXJVRWd5Z2t3bEdvbFBJVmxGQUpZQ3dZbGRlTEJsdk9BcHN3WEFmZ05kenVZMThaMEhwSTFhOWRIdko0MGFRRFZ1R21ZY3hSYXMrRnQ1MmRJQ3VvSU5DZWVOZEdwUzd4RkxndWVIaG0xY3JMSjBIRUhCVFpwRno0NWJGSG9vc2kweXoxL0RrNm9RUGhLZStkdmlOUVZocnp5VGtsTHNkZ1dGWW1NNVpVQTMvQT09IiwiaGFzaEFsZ29yaXRobSI6InNoYTI1NiJ9LCJzaWduYXR1cmUiOiI3N3lBWXpCdzR3UmVsN3A0T2lXcWxFZnNXeWFpQVA3VlJpVkhlNnlGWGN5WlArR1Q4UVg1cXdRdGk3TzAifQ==",
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
