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

// NewRancherSecretGetter returns a new rancherSecretGetter
func NewRancherSecretGetter(params *options) (SecretGetter, error) {
	client, err := newRancherClient()
	if err != nil {
		return &rancherSecretGetter{}, err
	}

	url := os.Getenv("CATTLE_URL")

	return &rancherSecretGetter{
		user:     os.Getenv("CATTLE_AGENT_ACCESS_KEY"),
		password: os.Getenv("CATTLE_AGENT_SECRET_KEY"),
		url:      strings.Replace(url, "v1", "v2-beta", 1),
		client:   client,
		token:    params.Token,
	}, nil
}

func (rsg rancherSecretGetter) GetSecrets(params *options) ([]secret, error) {
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
