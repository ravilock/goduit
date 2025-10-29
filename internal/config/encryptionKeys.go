package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
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

func LoadKeysFromEnv() error {
	privateKeyB64 := os.Getenv("JWT_PRIVATE_KEY_BASE64")
	if privateKeyB64 == "" {
		return nil
	}

	privateKeyContent, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return err
	}

	privatePem, _ := pem.Decode(privateKeyContent)
	PrivateKey, err = x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		return err
	}

	publicKeyB64 := os.Getenv("JWT_PUBLIC_KEY_BASE64")
	if publicKeyB64 == "" {
		return nil
	}

	publicKeyContent, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return err
	}

	return os.Setenv("PUBLIC_KEY", string(publicKeyContent))
}
