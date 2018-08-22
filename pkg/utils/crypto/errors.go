package crypto

import "errors"

var (
	InvalidPublicKeyError = errors.New("[crypto] invalid public key")
	InvalidPrivateKeyError = errors.New("[crypto] invalid private key")
	InvalidInputError = errors.New("[crypto] invalid input format")
	CipherTooShortError = errors.New("[crypto] cipher text too short")
)
