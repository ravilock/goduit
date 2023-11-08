package encryptionkeys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

var PrivateKey *rsa.PrivateKey

func LoadKeys(privateKeyFile, publicKeyFile *os.File) error {
	if err := readPrivateKey(privateKeyFile); err != nil {
		return err
	}
	return readPublicKey(publicKeyFile)
}

func readPrivateKey(privateKeyFile *os.File) error {
	stat, err := privateKeyFile.Stat()
	if err != nil {
		return err
	}
	privateKeyContent := make([]byte, stat.Size())

	if _, err := privateKeyFile.Read(privateKeyContent); err != nil {
		return err
	}

	privatePem, _ := pem.Decode(privateKeyContent)
	PrivateKey, err = x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		return err
	}
	return nil
}

func readPublicKey(publicKeyFile *os.File) error {
	stat, err := publicKeyFile.Stat()
	if err != nil {
		return err
	}
	publicKeyContent := make([]byte, stat.Size())

	if _, err := publicKeyFile.Read(publicKeyContent); err != nil {
		return err
	}
	return os.Setenv("PUBLIC_KEY", string(publicKeyContent))
}
