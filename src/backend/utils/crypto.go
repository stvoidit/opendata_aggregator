package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"hash/fnv"
	"io"
)

// Cryptographer - ...
type Cryptographer struct {
	key []byte
}

// NewCryptographer - ...
func NewCryptographer(key string) *Cryptographer {
	h := fnv.New128a()
	h.Write([]byte(key))
	return &Cryptographer{key: h.Sum(nil)}
}

// Decrypt - ...
func (e Cryptographer) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	var size = block.BlockSize()
	if len(data) < size {
		return nil, errors.New("text is too short")
	}
	iv := data[:size]
	data = data[size:]
	cipher.NewCTR(block, iv).XORKeyStream(data, data)
	return data, nil
}

// Encrypt - ...
func (e Cryptographer) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	var size = block.BlockSize()
	iv := genRandomIV(size)
	cipher.NewCTR(block, iv).XORKeyStream(data, data)
	return append(iv, data...), nil
}

// StringHMAC - ...
func (e Cryptographer) StringHMAC(value string) string {
	hash := hmac.New(sha512.New, e.key)
	io.WriteString(hash, value)
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func genRandomIV(size int) []byte {
	var iv = make([]byte, size)
	rand.Read(iv)
	return iv
}
