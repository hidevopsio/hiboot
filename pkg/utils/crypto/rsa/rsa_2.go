package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
)

var (
	publicKeyError  = errors.New("public key error")
	privateKeyError = errors.New("private key error!")
	parametersError = errors.New("parameters error")
)

func EncryptLongString(raw string, publicKey []byte) (result string, err error) {
	if len(raw) == 0 || len(publicKey) == 0 {
		return "", parametersError
	}
	public, err := getPublickKey(publicKey)
	if err != nil {
		return
	}

	partLen := public.N.BitLen()/8 - 11
	chunks := split([]byte(raw), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, _ := rsa.EncryptPKCS1v15(rand.Reader, public, chunk)
		buffer.Write(bytes)
	}
	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

func DecryptLongString(cipherText string, publicKey, privateKey []byte) (result string, err error) {
	if len(cipherText) == 0 || len(publicKey) == 0 || len(privateKey) == 0 {
		return "", parametersError
	}
	public, err := getPublickKey(publicKey)
	if err != nil {
		return
	}
	private, err := getPrivateKey(privateKey)
	if err != nil {
		return
	}
	partLen := public.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(cipherText)
	chunks := split([]byte(raw), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, private, chunk)
		buffer.Write(decrypted)
	}
	return buffer.String(), err
}

func getPublickKey(publicKey []byte) (*rsa.PublicKey, error) {
	if len(publicKey) == 0 {
		return nil, parametersError
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, publicKeyError
	}

	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	return pubInterface.(*rsa.PublicKey), nil
}

func getPrivateKey(priveteKey []byte) (*rsa.PrivateKey, error) {
	if len(priveteKey) == 0 {
		return nil, parametersError
	}
	block, _ := pem.Decode(priveteKey)
	if block == nil {
		return nil, privateKeyError
	}

	privInterface, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	return privInterface.(*rsa.PrivateKey), nil
}

func split(buf []byte, lim int) [][]byte {
	if buf == nil {
		return nil
	}
	if lim == 0 {
		return nil
	}
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

// 生成密钥对
func GenKeys(publicKeyWriter, privateKeyWriter io.Writer, keyLength int) error {
	if publicKeyWriter == nil || privateKeyWriter == nil {
		return parametersError
	}
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}
	derStream, _ := MarshalPKCS8PrivateKey(privateKey)

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	_ = pem.Encode(privateKeyWriter, block)

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, _ := x509.MarshalPKIXPublicKey(publicKey)
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	_ = pem.Encode(publicKeyWriter, block)

	return nil
}

func MarshalPKCS8PrivateKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, parametersError
	}
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(privateKey)
	k, _ := asn1.Marshal(info)
	return k, nil
}
