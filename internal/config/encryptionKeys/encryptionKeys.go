package encryptionkeys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
)

var PrivateKey *rsa.PrivateKey

func LoadPrivateKey(privateKey io.Reader) error {
	privateKeyContent, err := io.ReadAll(privateKey)
	if err != nil {
		return err
	}

	privatePem, _ := pem.Decode(privateKeyContent)
	PrivateKey, err = x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		return err
	}
	return nil
}

func LoadPublicKey(publicKeyFile io.Reader) error {
	publicKeyContent, err := io.ReadAll(publicKeyFile)
	if err != nil {
		return err
	}
	return os.Setenv("PUBLIC_KEY", string(publicKeyContent))
}
