package ports

type SignatureValidator interface {
	ValidateSignature(timestamp string, signature string, body []byte) bool
}
