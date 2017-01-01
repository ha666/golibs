package golibs

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
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
