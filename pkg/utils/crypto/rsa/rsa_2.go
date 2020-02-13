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

func EncryptLongString(raw string, publicKey []byte) (result string, err error) {
	public, err := getPublickKey(publicKey)
	if err != nil {
		return
	}

	partLen := public.N.BitLen()/8 - 11
	chunks := split([]byte(raw), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, public, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(bytes)
	}
	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

func DecryptLongString(cipherText string, publicKey, privateKey []byte) (result string, err error) {
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
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, private, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}
	return buffer.String(), err
}

func getPublickKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubInterface.(*rsa.PublicKey), nil
}

func getPrivateKey(priveteKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priveteKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privInterface.(*rsa.PrivateKey), nil
}

func split(buf []byte, lim int) [][]byte {
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
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}
	derStream := MarshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return err
	}
	return nil
}

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k
}
