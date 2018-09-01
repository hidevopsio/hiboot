package base64

import "encoding/base64"

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

// EncodeToString encode src string to base64 format string
func EncodeToString(src string) string {
	return coder.EncodeToString([]byte(src))
}

// Encode encode src bytes to base64 format bytes
func Encode(src []byte) (dst []byte) {
	dst = make([]byte, coder.EncodedLen(len(src)))
	coder.Encode(dst, src)
	return
}

// DecodeToString decode string from base64 string
func DecodeToString(src string) (retVal string, err error) {
	retBytes, err := coder.DecodeString(src)
	retVal = string(retBytes)
	return
}

// Decode decode base64 bytes
func Decode(src []byte) (dst []byte, err error) {
	size := coder.DecodedLen(len(src))
	buf := make([]byte, size)
	_, err = coder.Decode(buf, src)
	dst = buf[:size-1]
	return
}
