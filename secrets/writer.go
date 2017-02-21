package secrets

import (
	"github.com/rancher/secrets-api/pkg/aesutils"
)

// NewRSASecretFileWriter returns a SecretWriter implemenation to talk to Rancher
func NewRSASecretFileWriter(decryptor Decryptor) (SecretWriter, error) {
	return &rsaSecretFileWriter{
		decryptor: decryptor,
	}, nil
}

func (rsw rsaSecretFileWriter) Write(secrets []secret, dstDir string) error {
	for _, secret := range secrets {
		encData, err := getEncryptedData(secret.RewrapText)
		if err != nil {
			return err
		}

		aesKey, err := rsw.decryptor.Decrypt(encData.EncryptedKey.EncryptedText)
		if err != nil {
			return err
		}

		aesDecryptionKey := aesutils.NewAESKeyFromBytes(aesKey)

		clearText, err := aesutils.GetClearText(aesDecryptionKey, encData.EncryptedText)
		if err != nil {
			return err
		}

		err = secret.writeFile(dstDir, []byte(clearText))
		if err != nil {
			return err
		}
	}
	return nil
}
