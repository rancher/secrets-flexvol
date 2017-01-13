package secrets

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// SecretGetter gets the secrets froma remote source
type SecretGetter interface {
	GetSecrets(params map[string]interface{}) (*bulkSecret, error)
}

type rancherSecretGetter struct {
	url    string
	client *http.Client
}

// NewRancherSecretGetter returns a new rancherSecretGetter
func NewRancherSecretGetter(params map[string]interface{}) (SecretGetter, error) {
	client, err := newRancherClient()
	if err != nil {
		return &rancherSecretGetter{}, err
	}

	authPair := os.Getenv("CATTLE_AGENT_ACCESS_KEY") + ":" + os.Getenv("CATTLE_AGENT_SECRET_KEY")
	url := os.Getenv("CATTLE_URL")
	urlSplit := strings.SplitN(url, "//", 2)

	if len(urlSplit) == 2 {
		url = urlSplit[0] + "//" + authPair + "@" + urlSplit[1]
	}

	return &rancherSecretGetter{
		url:    url,
		client: client,
	}, nil
}

func (rsg rancherSecretGetter) GetSecrets(params map[string]interface{}) (*bulkSecret, error) {
	reqURL := rsg.url
	returnSecrets := &bulkSecret{}

	token, ok := params["io.rancher.secrets.token"]
	if !ok {
		return returnSecrets, errors.New("No token found")
	}

	logrus.Infof("REQUESTING SECRETS FROM: %s", reqURL)
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return returnSecrets, err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(tokenBytes))
	if err != nil {
		return returnSecrets, err
	}

	req.Header.Add("Content-Type", "application/x-api-secrets-token")

	resp, err := rsg.client.Do(req)
	if err != nil {
		return returnSecrets, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return returnSecrets, err
	}
	logrus.Infof("rancher-resp: %#v", body)

	return returnSecrets, json.Unmarshal(body, returnSecrets)
}

func newRancherClient() (*http.Client, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return client, nil
}
