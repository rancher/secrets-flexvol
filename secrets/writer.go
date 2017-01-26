package secrets

// NewRSASecretFileWriter returns a SecretWriter implemenation to talk to Rancher
func NewRSASecretFileWriter(decryptor Decryptor) (SecretWriter, error) {
	return &rsaSecretFileWriter{
		decryptor: decryptor,
	}, nil
}

func (rsw rsaSecretFileWriter) Write(secrets []secret, dstDir string) error {
	for _, secret := range secrets {
		clearText, err := rsw.decryptor.Decrypt(secret.RewrapText)
		if err != nil {
			return err
		}

		err = secret.writeFile(dstDir, clearText)
		if err != nil {
			return err
		}
	}
	return nil
}
