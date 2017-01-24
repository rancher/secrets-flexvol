package secrets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// SecretGetter gets the secrets froma remote source
type SecretGetter interface {
	GetSecrets(params map[string]interface{}) ([]secret, error)
}

type rancherSecretGetter struct {
	user     string
	password string
	url      string
	client   *http.Client
	token    *rancherToken
}

type rancherToken struct {
	Value []byte `json:"value"`
}

// NewRancherSecretGetter returns a new rancherSecretGetter
func NewRancherSecretGetter(params map[string]interface{}) (SecretGetter, error) {
	rToken := &rancherToken{}

	client, err := newRancherClient()
	if err != nil {
		return &rancherSecretGetter{}, err
	}

	url := os.Getenv("CATTLE_URL")

	if rawToken, ok := params["io.rancher.secrets.token"].(string); ok {
		rToken.Value = []byte(strings.Replace(rawToken, "\\", "", -1))
	}

	return &rancherSecretGetter{
		user:     os.Getenv("CATTLE_AGENT_ACCESS_KEY"),
		password: os.Getenv("CATTLE_AGENT_SECRET_KEY"),
		url:      strings.Replace(url, "v1", "v2-beta", 1),
		client:   client,
		token:    rToken,
	}, nil
}

func (rsg rancherSecretGetter) GetSecrets(params map[string]interface{}) ([]secret, error) {
	reqURL := rsg.url + "/secrets"
	returnSecrets := []secret{}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(rsg.token.Value))
	if err != nil {
		return returnSecrets, err
	}

	req.Header.Add("Content-Type", "application/x-api-secrets-token")
	req.SetBasicAuth(rsg.user, rsg.password)

	resp, err := rsg.client.Do(req)
	if err != nil {
		logrus.Errorf("Response code: %s", resp.Status)
		return returnSecrets, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return returnSecrets, fmt.Errorf("Unsuccessful request: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return returnSecrets, err
	}

	err = json.Unmarshal(body, &returnSecrets)

	return returnSecrets, err
}

func newRancherClient() (*http.Client, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return client, nil
}
