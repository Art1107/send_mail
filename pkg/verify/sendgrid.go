package verify

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
)

type ecdsaSignature struct {
	R *big.Int
	S *big.Int
}

func VerifySignature(payload []byte, signature string, timestamp string, publicKey string) (bool, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false, fmt.Errorf("failed to parse PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %v", err)
	}

	ecdsaKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("key is not an ECDSA public key")
	}

	decodedSig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	var sig ecdsaSignature
	if _, err := asn1.Unmarshal(decodedSig, &sig); err != nil {
		return false, fmt.Errorf("failed to unmarshal signature: %v", err)
	}
	h := sha256.New()
	h.Write([]byte(timestamp))
	h.Write(payload)
	hash := h.Sum(nil)

	details := fmt.Sprintf("Payload length: %d, Timestamp: %s", len(payload), timestamp)
	fmt.Printf("Verification Details: %s\n", details)

	isValid := ecdsa.Verify(ecdsaKey, hash, sig.R, sig.S)

	if isValid {
		fmt.Printf("Signature verification successful\n")
	} else {
		fmt.Printf("Signature verification failed\n")
	}

	return isValid, nil
}
