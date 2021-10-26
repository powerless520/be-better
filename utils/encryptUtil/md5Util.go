package encryptUtil

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(msg string) string {
	h := md5.New()
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}