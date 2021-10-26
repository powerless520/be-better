package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func HMACEncode(msg string, key string) string {
	hmacCtx := hmac.New(sha256.New, []byte(key))
	hmacCtx.Write([]byte(msg))
	cipherStr := hmacCtx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func HMACSha1Encode(msg string, key string) string {
	hmacCtx := hmac.New(sha1.New, []byte(key))
	hmacCtx.Write([]byte(msg))
	cipherStr := hmacCtx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func HMACSha1Base64Encode(msg string, key string) string {
	hmacCtx := hmac.New(sha1.New, []byte(key))
	hmacCtx.Write([]byte(msg))
	cipherStr := hmacCtx.Sum(nil)
	return base64.StdEncoding.EncodeToString(cipherStr)
}
