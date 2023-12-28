package encryptionkeys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

const privateKeyFileName = "jwtRS256.key"
const publicKeyFileName = "jwtRS256.key.pub"

var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

func LoadKeys() error {
	if err := readPrivateKey(); err != nil {
		return err
	}
	return readPublicKey()
}

func readPrivateKey() error {
	privateKeyContent, err := os.ReadFile(privateKeyFileName)
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

func readPublicKey() error {
	publicKeyContent, err := os.ReadFile(publicKeyFileName)
	if err != nil {
		return err
	}

	publicPem, _ := pem.Decode(publicKeyContent)
	PublicKey, err = x509.ParsePKCS1PublicKey(publicPem.Bytes)
	if err != nil {
		return err
	}
	return nil
}
