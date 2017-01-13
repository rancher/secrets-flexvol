package secrets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// Decryptor handles decrypting messages
type Decryptor interface {
	Decrypt(cipherText string) ([]byte, error)
}

type rsaDecryptor struct {
	privateKeyPath string
}

// NewRSADecryptor returns an RSA decryptor
func NewRSADecryptor(privateKeyPath string) (Decryptor, error) {
	return rsaDecryptor{
		privateKeyPath: privateKeyPath,
	}, nil
}

// Decrypt implments the decryptor interface
func (r rsaDecryptor) Decrypt(cipherText string) ([]byte, error) {
	key, err := loadPrivateKeyFromFile(r.privateKeyPath)
	if err != nil {
		return []byte{}, err
	}

	return rsaDecrypt(key, cipherText)
}

func loadPrivateKeyFromFile(keyPath string) (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("Could not decode private key. Is it PEM format?")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func loadPrivateKeyFromString(keyString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyString))
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func rsaDecrypt(priv *rsa.PrivateKey, cipherText string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return []byte{}, err
	}

	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, data, []byte(""))
}
