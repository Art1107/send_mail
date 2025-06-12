package adapters

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"log"
)

type SendgridValidator struct {
	publicKey string
}

func NewSendgridValidator(publicKey string) *SendgridValidator {
	return &SendgridValidator{publicKey: publicKey}
}

func (v *SendgridValidator) ValidateSignature(timestamp string, signature string, body []byte) bool {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(v.publicKey)
	if err != nil {
		log.Printf("Error decoding public key: %v", err)
		return false
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		log.Printf("Error parsing public key: %v", err)
		return false
	}

	ecdsaKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Printf("Public key is not an ECDSA key")
		return false
	}

	message := append([]byte(timestamp), body...)

	hash := sha256.Sum256(message)

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		log.Printf("Error decoding signature: %v", err)
		return false
	}

	return ecdsa.VerifyASN1(ecdsaKey, hash[:], signatureBytes)
}
