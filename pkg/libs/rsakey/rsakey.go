package rsakey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// GenerateRSAKeys creates a pair of private and public keys for a client.
func GenerateRSAKeys() (string, string, error) {
	/* Shamelessly borrowed and adapted from some golang-samples */
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	if err := priv.Validate(); err != nil {
		errStr := fmt.Errorf("RSA key validation failed: %s", err)
		return "", "", errStr
	}
	privDer := x509.MarshalPKCS1PrivateKey(priv)
	/* For some reason chef doesn't label the keys RSA PRIVATE/PUBLIC KEY */
	privBlk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDer,
	}
	privPem := string(pem.EncodeToMemory(&privBlk))
	pub := priv.PublicKey
	pubDer, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		errStr := fmt.Errorf("Failed to get der format for public key: %s", err)
		return "", "", errStr
	}
	pubBlk := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDer,
	}
	pubPem := string(pem.EncodeToMemory(&pubBlk))
	return privPem, pubPem, nil
}
