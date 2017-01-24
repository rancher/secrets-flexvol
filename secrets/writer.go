package secrets

// SecretWriter implements the Writer interface
type SecretWriter interface {
	Write(secrets []secret, dst string) error
}

type rsaSecretFileWriter struct {
	decryptor Decryptor
	params    map[string]interface{}
}

// NewRSASecretFileWriter returns a SecretWriter implemenation to talk to Rancher
func NewRSASecretFileWriter(decryptor Decryptor, params map[string]interface{}) (SecretWriter, error) {
	return &rsaSecretFileWriter{
		decryptor: decryptor,
		params:    params,
	}, nil
}

func (rsw rsaSecretFileWriter) Write(secrets []secret, dstDir string) error {
	for _, secret := range secrets {
		plainText, err := rsw.decryptor.Decrypt(secret.RewrapText)
		if err != nil {
			return err
		}

		err = secret.writeFile(dstDir, plainText)
		if err != nil {
			return err
		}
	}
	return nil
}
