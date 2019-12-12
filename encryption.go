package golibs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func Md5(s string) string {
	h := md5.New()
	h.Write(StringToSliceByte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func UrlEncode(source string) string {
	return url.QueryEscape(source)
}

func UrlDecode(source string) (string, error) {
	return url.QueryUnescape(source)
}

func Base64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func UnBase64(source string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(source)
}

func Sha1(message []byte) []byte {
	h := sha1.New()
	h.Write(message)
	return h.Sum(nil)
}

func HmacMd5(message, secret []byte) []byte {
	h := hmac.New(md5.New, secret)
	h.Write(message)
	return h.Sum(nil)
}

func HmacSha1(message, secret []byte) []byte {
	h := hmac.New(sha1.New, secret)
	h.Write(message)
	return h.Sum(nil)
}

func HmacSha256(message, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(message)
	return h.Sum(nil)
}

func HmacSha512(message, secret []byte) []byte {
	h := hmac.New(sha512.New, secret)
	h.Write(message)
	return h.Sum(nil)
}

func AesEncrypt(plaintext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func AesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7Unpadding(origData)
	return origData, nil
}

func DesEncrypt(origData, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS7Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7Unpadding(origData)
	return origData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func BcryptEncrypt(message []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(message, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func CheckBcrypt(hashBytes, s []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashBytes, s); err != nil {
		return false, err
	} else {
		return true, nil
	}
}
