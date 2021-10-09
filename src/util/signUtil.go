package util

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

func RsaVerifyWithSha1(ori_data, sign_data, pubKey string) error {
	//sign, err := base64.StdEncoding.DecodeString(sign_data)
	//if err != nil {
	//	return err
	//}
	// 转16进制字节数组
	sign, err := hex.DecodeString(sign_data)
	if err != nil {
		return err
	}

	block,_ := pem.Decode([]byte(pubKey))
	if block == nil {
		return errors.New("blockNil")
	}

	//pubKeyBy, err := base64.StdEncoding.DecodeString(pubKey)
	//if err != nil {
	//	return err
	//}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(ori_data))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}
