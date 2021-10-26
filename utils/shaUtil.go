package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func SHA256Encode(msg string) string {
	md5Ctx := sha256.New()
	md5Ctx.Write([]byte(msg))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func SHA1Encode(msg string) string {
	md5Ctx := sha1.New()
	md5Ctx.Write([]byte(msg))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
