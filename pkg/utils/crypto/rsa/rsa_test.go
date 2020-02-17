package rsa

import (
	"bytes"
	"crypto/rsa"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/crypto"
)

var invalidPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MCsCAQACBQDJaQRdAgMBAAECBEPxlU0CAwDt4wIDANi/AgJC0QICBc0CAkg4
-----END RSA PRIVATE KEY-----
`)

var invalidPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MCAwDQYJKoZIhvcNAQEBBQADDwAwDAIFAMlpBF0CAwEAAQ==
-----END PUBLIC KEY-----
`)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestRsa(t *testing.T) {
	src := []byte("hello")
	data, _ := Encrypt([]byte(src))
	decrypted, err := Decrypt(data)
	assert.Equal(t, nil, err)
	//log.Debugf("encrypted: %v, decrypted: %v, org: %v", string(data), string(decrypted), string(src))
	assert.Equal(t, src, decrypted)
}

func TestRsaBase64(t *testing.T) {
	src := []byte("hello")
	data, _ := EncryptBase64([]byte(src))
	decrypted, err := DecryptBase64(data)
	assert.Equal(t, nil, err)
	//log.Debugf("encrypted: %v, decrypted: %v, org: %v", string(data), string(decrypted), string(src))
	assert.Equal(t, src, decrypted)
}

func TestExeptions(t *testing.T) {
	t.Run("should report error with invalid public key", func(t *testing.T) {
		src := []byte("hello")
		_, err := Encrypt([]byte(src), []byte("invalid-key"))
		assert.Equal(t, crypto.ErrInvalidPublicKey, err)
	})

	t.Run("should report error with invalid public key", func(t *testing.T) {
		src := []byte("hello")
		_, err := Encrypt([]byte(src), invalidPrivateKey)
		assert.Contains(t, err.Error(), "asn1: structure error")
	})

	t.Run("should report error with invalid public key", func(t *testing.T) {
		src := []byte("hello")
		data, _ := Encrypt([]byte(src))
		_, err := Decrypt(data, invalidPublicKey)
		assert.Contains(t, err.Error(), "tags don't match")
	})

	t.Run("should report error with invalid private key", func(t *testing.T) {
		src := []byte("hello")
		data, _ := Encrypt([]byte(src))
		_, err := Decrypt(data, []byte("invalid-key"))
		assert.Equal(t, crypto.ErrInvalidPrivateKey, err)
	})

	t.Run("should report error with invalid base64 string", func(t *testing.T) {
		src := []byte("invalid-base64")
		_, err := DecryptBase64([]byte(src))
		assert.Equal(t, crypto.ErrInvalidInput, err)
	})
}

var testPriveteKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEuwIBADALBgkqhkiG9w0BAQEEggSnMIIEowIBAAKCAQEArn3nYQSwTsm9hjDf
7LiM6yM2ToFh8Guxxsld6FPQ2j7CBei4scO6YjSS4TvMeeufznQS9bzWvekYdirT
77eFV5MjfXpiylHFHNS5gcZihThvPJrOVagyyAV11Abi+EY9/1aRrTGF6JjX+jqF
2/jwhwGkmhefvoxHRCAjkemi9QO8omxh/AU4c3EjyWMVjaCgngtjvkiBeMU9Wsag
ZiYnEez01/fP4pq/n1pWBc7QtmH20Z5NOj9LPC2+8oYGpke9lG/IdfGuqicJJNC/
0vZEJoK4D3+BiokvR20OC9c43ffznfbhGV9lblyOS30kIPcyAlknJrWVkKcNIhIL
PMrh4QIDAQABAoIBADCgtdK708ahQkgbZsw5wkvlTEUkmX6/BJQ5mgodEZ9AziGH
cbFYsqCbtjM+zwVLPQX0IzSIo+/Y/hAwb0/m/Soiv0lAyjdIAn6+adRYzSwDRjzF
h6snbL+BhgzIvogiSzTVk1OI8aCYt9fsZ1GeVqnJM24eF06rGVFLA56uVdOh+Sl/
YefEnAYHTsbsF3hEyk2WsnqlgxqRIGiK5TCpp1uz8UZvt73Qkc5L3gVTftkuMFk8
NpLlQI/FwME0+K//lgFBpvBW6cvoe6TwsL67q6yxYQ4Nma22lx2wUeSuaHjCmRE0
IbzDpc2IlIqr8T0WOTetaiUcSfz+37Z4ZqB+4gECgYEA1R/nRXh8UcPiXYQNeJOG
I25rqYPOMMpRHsWxFYkRrXOv9p6kLXpZ4cB15DJoMZrETVVU0BHkvkWZ1aZatSO4
OuyMc4PIYRdQoM/LH+jeLuTSY/yu1b6jrYxxOYQSPxcqBCLXnRzOl6CFxybQtMAx
gwXlPaZfH8o5cvKoLu10nMkCgYEA0Zhg2DyTRr81rGSUUXNCfAgG45Yz7VSTXOeu
E0efaGL8VaeKKBrD1/GKATo0qetT+BlSWZ+4zmDnVLxVkOhfqFIcZ6Tug+sqoDXm
e+2CIJOzDRnl+uqBTd1qn6VBmHTneemAG35z1mv4Khx7J1/FQj8UPKsxzh43YD/l
P2xNYFkCgYEAlc8GLxQBNyxdCuUO3wm7yV4NuatXnX4pRVsrHfsqfOUL9GwQ9ZLC
aWhytgQkr3GduMpZgqSBSKn992sm6ZsBHhI2q+AfUvgjidZmbriurQHVTclJUB/g
R9anpAlNFiH/O8cODnc4VObWAmYrYFKUuwfC2vH+fYcVmNIvHEV3qdkCgYByMtAx
gW/NUElyUKrvZhmHcuguAJzyZu6T5DfYkWGtgqFyGgMQruSeOCC1Yn1nR61MtJ9F
7dzHtczVQnhsp+/WykZnwlmizvM+r5+RTmtkTJV2QfIosLUbM9Twfx4qbyfgKPWA
BXogDlv8td/0KB5WZgAkvjI42AXcD3RdBilyoQKBgBOVhHhVYMdKY7fibpXdL/W8
7+Boy8mRuRxY06Df6d1R2LuiLFCpKVknP+oWx90dFiA0cXKNw57cvacImZxt3V1+
+beiFhbYGVnSkvDOjsFF3zZfK00df+4mT211DSYegqk/iYUBA+Z5QBxVWO+vApKX
wvloIFUua5y1py6nksjU
-----END PRIVATE KEY-----
`)

var testPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArn3nYQSwTsm9hjDf7LiM
6yM2ToFh8Guxxsld6FPQ2j7CBei4scO6YjSS4TvMeeufznQS9bzWvekYdirT77eF
V5MjfXpiylHFHNS5gcZihThvPJrOVagyyAV11Abi+EY9/1aRrTGF6JjX+jqF2/jw
hwGkmhefvoxHRCAjkemi9QO8omxh/AU4c3EjyWMVjaCgngtjvkiBeMU9WsagZiYn
Eez01/fP4pq/n1pWBc7QtmH20Z5NOj9LPC2+8oYGpke9lG/IdfGuqicJJNC/0vZE
JoK4D3+BiokvR20OC9c43ffznfbhGV9lblyOS30kIPcyAlknJrWVkKcNIhILPMrh
4QIDAQAB
-----END PUBLIC KEY-----
`)

func TestEncryptDecryptLongString(t *testing.T) {
	random2048Str := "oqrjv5SaMqEjgbucY8bLI8ugwvzCTTKdkcHRzOmRQRSRUYecboVByaijWZyem5pnl7K4a8dnvUlEj4xU2TQ460wgxwIYHoXoe62m90ae0K4DPhJI39ajvbHf0F5vaf6cmyT5YSkansmYzjNyOfSmUcrvDhj2DGgik4psGrFEiJ1qnwUuHZszqG0yuWDkg0YnWyYdnigFJEU2KuJNCm7z65LWjOvX1Jz2yMNKfuj8gzeUPPfxaJ8auG1xQF54ORjUPvzCRuenMyRfFHvnaGmt6HtmNs6Ex7WodFptEcs3ptUkXlMNKfkKYe8719WRadQQcFqsW13puSbWhoopaaEHTmW5GQcoLK9dNrzPhY4AgSM4lr8pvgGf6exv1RanLWg3zCSeC8dfuMcmcHzLwdiTLJHCSY2KW6pdGvM7DYyXsUL2hvKlAGvCuBKGAxMgqit2uYZykNnJ5UnrF4fIb9QnbtZyPtAe2WSpoD8cKAbrm4MuTbSSw7daKWHlT7dL2CnaY6TPG4oHn0NCTxvF5G3gJKEi1oUI0vdKB84WkBC6Rc8e5k6i2c3EnyXg6XI1uHbuUNj4DyQr2oYoy5RPlHWSXUCAK04i0HT2iCH9bKKL7DgWnSDkhq13S2kDjPmjTbeVjsA2cbzse4p49cihxT8vAeMM0LunpJRiyhA4DUwdbKltMbpMT78RzWKtUdqa5iQORspD7iGwaMouyhssEJDIVYCaz8n1CX8r6NKMyXx8Bu8mXtLpqPGGDI2NbCPJXMtVCv9NfHYRJg1Jolb7HnSUGNL1bzp9WFzjWm7VUu4VmGOBHOkWJtmX3x9vGwuZvBOIiPiTKs9txzQsNkO0XXmWcaphqXZkssJ6WqXkTNHMgaxvNQEb2EyT2uYuwcD3upv0KODdfu4e2HQroWvkaesARghACS1HawfAXqO8N56Xx1kXxW5TU7dZJICDrBqwzKDra5A7dHcRMkqqAghx6jbv0v8ZbqyR3JDDoRfwzZlMOVPk7qGjwDj2Lyp2BAvCCrQDY6xbVyZgqtCtUUjWRKaAPbJpQeapCVQW53JTBf2Jmg36182SVkNdcgrreM40MFtvMqtZG9wJ9SMEgwU5G4mye98gybFQYShsaLvqdx8yW9wFgvuuvspbIQV3819gji9eeiZiPSzDJZYI33hi4AFj1x3YmoMaSQcATplgJuaRADG7l72cN6vG4FHcl9vGPAZZqNtjRPmbh2cSG5QEzaoCXvA3RdnBvkQuFiQzrxvVbABfGD7fGnRhZ88l7iy2qtl1kaoU28N879k0GMGlgga4LccktZGDrH9kGfN35cYli99SMIZT9hvY24oN1clDR1KNin0YYrGr2uhsaKXgwkFJmZBHEvRrRTuIS07ZJ1Eck10OjK9IxTzQ1D1fYAbqclQGO4ZZBu7PHOgaRYomsJUwmJR76T7lxOq3g1WkDpXoQn72fZtyFV9ga8CYF4ozE15CFvS6TW3fgWiR3PUUzujbuE8xofhNVhAnyDkBHZKiMicpxNpBPv3tOhwk3dH9lwQSvLRISja3e9aCZkyFEkmcxxQWKv5OmZN47cjSC6Ifwq33muBq9TfLQPaGcdaSsLzCR98HZOgDizvXhRgZjRiR2WMwwMkFw0Y8pS9FeHPOCWX5V8RvtH3mHFN9b9P2OInTdrsIrjoQwlFbEKa8bfbHryoeZ26BmPNA1HNvOdW84no7XZTP3PkbwyOIftTgBiSUZi4kr8COKl6eJx2XcagaBhj9KlAStQUfKMLgu0lWSkRBlEOQqi63kezdxgIzQcdJiEzqXvcfBgpwXK76nYfSl7aufcTIPisHUQiuv0xniALrTJwaq6snCGTOdFuecF6nqt11Ih9B2efcmK8phjtyyL8ML1iIN8cnPPpRzUGUdChGmqjtBMeXOKezjMPZRvRXmwCey3s0gWEQkqSr18rrUWeXqNBXYxbjMkTfVOo4S6psiJy3pIjPboLlOTSwf6twNuBUPGEfEErwZgs5peO8DTvOtk81x2BP1Yh8fjTvBFOCqeN3LhtpViwUkQXQ8zj7"
	encryptResult, err := EncryptLongString(random2048Str, testPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	decryptResult, err := DecryptLongString(encryptResult, testPublicKey, testPriveteKey)
	if err != nil {
		t.Fatal(err)
	}
	if decryptResult != random2048Str {
		t.Fatal("They should be equal")
	}
}

func Test_split(t *testing.T) {
	type args struct {
		buf []byte
		lim int
	}
	tests := []struct {
		name string
		args args
		want [][]byte
	}{
		{
			"buf is nil",
			args{nil, 0},
			nil,
		},
		{
			"lim is 0",
			args{[]byte("123"), 0},
			nil,
		},
		{
			"split buf to buf size",
			args{[]byte("abc"), 1},
			[][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
			},
		},
		{
			"split buf to 2,but buf is 2",
			args{[]byte("ab"), 2},
			[][]byte{
				[]byte("ab"),
			},
		},
		{
			"split buf to 2,but buf len is 1",
			args{[]byte("a"), 2},
			[][]byte{
				[]byte("a"),
			},
		},
		{
			"split buf to 2,but buf len is 3",
			args{[]byte("abc"), 2},
			[][]byte{
				[]byte("ab"),
				[]byte("c"),
			},
		},
		{
			"split buf to 3,but buf len is 2",
			args{[]byte("ab"), 3},
			[][]byte{
				[]byte("ab"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.args.buf, tt.args.lim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryptLongString(t *testing.T) {
	type args struct {
		raw       string
		publicKey []byte
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{
			"raw len is 0",
			args{
				raw: "",
			},
			"",
			true,
		},
		{
			"publciKey len is 0",
			args{
				"a",
				[]byte{},
			},
			"",
			true,
		},
		{
			"publciKey err",
			args{
				"a",
				[]byte("abc"),
			},
			"",
			true,
		},
		//The results change every time
		{
			"publciKey right",
			args{
				"a",
				testPublicKey,
			},
			"abc",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncryptLongString(tt.args.raw, tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptLongString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if gotResult != tt.wantResult {
			//	t.Errorf("EncryptLongString() = %v, want %v", gotResult, tt.wantResult)
			//}
		})
	}
}

func TestDecryptLongString(t *testing.T) {
	var cipherText = "Nf5ugWi1uAQ_nuLUMuiSamyhYXkOgVaXraPk7U9U4oN6FzRJrXNzZPMEyPwSllKE-rLPWPItxpuVeDQWuU1VqYZ33G0DEUNcmVgWJbMug1njbFZngsQP0e3TIVjiAf5S5_O5F3ZbAOrZ6XN9eR5_ynCIcFLbpwFuPJRnLUYjA06ulYcQvR32hpYmcrx3Cj8jakcPwMDR_s1_CxaT-HOQFETZ4lpBRS4UGl-fez-oVM3PbEVf3AYFh4bN1KqW0EMu0glNoPd462Y7M6CZPch1DaZRml5lcq6H8NlQs092dG-574EWR3Dbbm7PHUxehfI91je4lnbFb2jkV66k5m9wzQ"
	type args struct {
		cipherText string
		publicKey  []byte
		privateKey []byte
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{
			"cipherText is empty",
			args{
				cipherText: "",
			},
			"",
			true,
		},
		{
			"publicKey len is 0",
			args{
				cipherText: "a",
				publicKey:  []byte{},
			},
			"",
			true,
		},
		{
			"privateKye len is 0",
			args{
				cipherText: "a",
				publicKey:  []byte("a"),
				privateKey: []byte{},
			},
			"",
			true,
		},
		{
			"publicKey error",
			args{
				cipherText: cipherText,
				publicKey:  []byte("a"),
				privateKey: []byte("a"),
			},
			"",
			true,
		},
		{
			"publicKey right",
			args{
				cipherText: cipherText,
				publicKey:  testPublicKey,
				privateKey: []byte("a"),
			},
			"",
			true,
		},
		{
			"publicKey right,privateKey right",
			args{
				cipherText: cipherText,
				publicKey:  testPublicKey,
				privateKey: testPriveteKey,
			},
			"a",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := DecryptLongString(tt.args.cipherText, tt.args.publicKey, tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptLongString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("DecryptLongString() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_getPublickKey(t *testing.T) {
	type args struct {
		publicKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *rsa.PublicKey
		wantErr bool
	}{
		{
			"publicKey len is 0",
			args{
				[]byte{},
			},
			nil,
			true,
		},
		{
			"publickey error",
			args{
				[]byte("abc"),
			},
			nil,
			true,
		},
		{
			"publicKey right",
			args{
				testPublicKey,
			},
			&rsa.PublicKey{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getPublickKey(tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPublickKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("getPublickKey() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Test_getPrivateKey(t *testing.T) {
	type args struct {
		priveteKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *rsa.PrivateKey
		wantErr bool
	}{
		{
			"privateKey len is 0",
			args{
				[]byte{},
			},
			nil,
			true,
		},
		{
			"privateKey error",
			args{
				[]byte("abc"),
			},
			nil,
			true,
		},
		{
			"privateKey right",
			args{
				testPriveteKey,
			},
			&rsa.PrivateKey{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getPrivateKey(tt.args.priveteKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("getPrivateKey() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestGenKeys(t *testing.T) {
	type args struct {
		pubW      io.Writer
		priW      io.Writer
		keyLength int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"publicKeyWriter is nil",
			args{
				nil,
				nil,
				0,
			},
			true,
		},
		{
			"privateKeyWriter is nil",
			args{
				&bytes.Buffer{},
				nil,
				0,
			},
			true,
		},
		{
			"key is 0",
			args{
				&bytes.Buffer{},
				&bytes.Buffer{},
				0,
			},
			true,
		},
		{
			"key is 2",
			args{
				&bytes.Buffer{},
				&bytes.Buffer{},
				2,
			},
			true,
		},
		{
			"key is 2048",
			args{
				&bytes.Buffer{},
				&bytes.Buffer{},
				2048,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//publicKeyWriter := &bytes.Buffer{}
			//privateKeyWriter := &bytes.Buffer{}
			if err := GenKeys(tt.args.pubW, tt.args.priW, tt.args.keyLength); (err != nil) != tt.wantErr {
				t.Errorf("GenKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if gotPublicKeyWriter := publicKeyWriter.String(); gotPublicKeyWriter != tt.wantPublicKeyWriter {
			//	t.Errorf("GenKeys() = %v, want %v", gotPublicKeyWriter, tt.wantPublicKeyWriter)
			//}
			//if gotPrivateKeyWriter := privateKeyWriter.String(); gotPrivateKeyWriter != tt.wantPrivateKeyWriter {
			//	t.Errorf("GenKeys() = %v, want %v", gotPrivateKeyWriter, tt.wantPrivateKeyWriter)
			//}
		})
	}
}

func TestToolGenEncryptStr(t *testing.T) {
	result, err := EncryptLongString("a", testPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestMarshalPKCS8PrivateKey(t *testing.T) {
	pk, _ := getPrivateKey(testPriveteKey)
	var errPk = *pk
	errPk.E = 1000
	type args struct {
		key *rsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"privateKey is nil",
			args{
				key: nil,
			},
			nil,
			true,
		},
		{
			"privateKey is error",
			args{
				key: &errPk,
			},
			[]byte{},
			false,
		},
		{
			"privateKey is right",
			args{
				key: pk,
			},
			[]byte{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MarshalPKCS8PrivateKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalPKCS8PrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("MarshalPKCS8PrivateKey() = %v, want %v", got, tt.want)
			//}
		})
	}
}
