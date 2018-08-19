package crypto

import "errors"

var (
	InvalidPublicKeyError = errors.New("invalid public key")
	InvalidPrivateKeyError = errors.New("invalid private key")
)
