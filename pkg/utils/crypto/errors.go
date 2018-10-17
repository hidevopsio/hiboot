package crypto

import "errors"

var (
	// ErrInvalidPublicKey invalid public key
	ErrInvalidPublicKey = errors.New("[crypto] invalid public key")
	// ErrInvalidPrivateKey invalid private key
	ErrInvalidPrivateKey = errors.New("[crypto] invalid private key")
	// ErrInvalidInput invalid input format
	ErrInvalidInput = errors.New("[crypto] invalid input format")
	// ErrCipherTooShort cipher text too short
	ErrCipherTooShort = errors.New("[crypto] cipher text too short")
)
