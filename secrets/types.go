package secrets

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
)

const (
	// DefaultMode is readable only by user
	DefaultMode = "0440"
	//DefaultUID is Root
	DefaultUID = "0"
	//DefaultGID is Roots default group
	DefaultGID = "0"
)

type secret struct {
	Name       string `json:"name"`
	UID        string `json:"uid"`
	GID        string `json:"gid"`
	Mode       string `json:"mode"`
	RewrapText string `json:"rewrapText"`
}

type options struct {
	Token   *secretToken `json:"io.rancher.secrets.token,omitempty"`
	Rancher bool         `json:"rancher,string,omitempty"`
	Device  string       `json:"device,omitempty"`
	Name    string       `json:"name,omitempty"`
}

type secretToken struct {
	Value []byte `json:"value,omitempty"`
}

// SecretGetter gets the secrets froma remote source
type SecretGetter interface {
	GetSecrets(params *options) ([]secret, error)
}

type rancherSecretGetter struct {
	user     string
	password string
	url      string
	client   *http.Client
	token    *secretToken
}

// SecretWriter implements the Writer interface
type SecretWriter interface {
	Write(secrets []secret, dst string) error
}

type rsaSecretFileWriter struct {
	decryptor Decryptor
}

func newOptions(params map[string]interface{}) (*options, error) {
	option := &options{}
	token := &secretToken{}

	//clean the token...
	tkn, ok := params["io.rancher.secrets.token"].(string)
	if !ok {
		return option, errors.New("No token passed")
	}
	delete(params, "io.rancher.secrets.token")
	token.Value, _ = clean([]byte(tkn))

	paramBytes, err := json.Marshal(params)
	if err != nil {
		return option, err
	}

	cParamBytes, _ := clean(paramBytes)

	err = json.Unmarshal(cParamBytes, option)
	if err != nil {
		logrus.Error(err)
		return option, err
	}
	option.Token = token

	return option, nil
}

func clean(val []byte) ([]byte, error) {
	stringToken := string(val)
	return []byte(strings.Replace(stringToken, "\\", "", -1)), nil
}
