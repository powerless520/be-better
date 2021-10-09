package util

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func RsaVerifyWithPubkey(origin_str, rsa_pubkey, sign_received string) error {
	sign_received_base64decode, err := base64.StdEncoding.DecodeString(sign_received)
	if err != nil {
		return err
	}
	block, _ := pem.Decode([]byte(rsa_pubkey))
	if block == nil {
		return errors.New("PublicKeyError")
	}
	//
	//public, err := base64.StdEncoding.DecodeString(rsa_pubkey)
	//if err != nil{
	//	return err
	//}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(origin_str))
	err = rsa.VerifyPKCS1v15(pubInterface.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign_received_base64decode)
	if err != nil {
		return err
	}
	return nil
}

// RsaVerifyWithPubkeyMD5 rsa解密md5校验签名
func RsaVerifyWithPubkeyMD5(originStr, rsaPubkey, signReceived string) error {
	signReceivedBase64decode, err := base64.StdEncoding.DecodeString(signReceived)
	if err != nil {
		return err
	}
	block, _ := pem.Decode([]byte(rsaPubkey))
	if block == nil {
		return errors.New("PublicKeyError")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hash := md5.New()
	hash.Write([]byte(originStr))
	err = rsa.VerifyPKCS1v15(pubInterface.(*rsa.PublicKey), crypto.MD5, hash.Sum(nil), signReceivedBase64decode)
	if err != nil {
		return err
	}
	return nil
}

func RsaPubKeyToSign(data, pub_key string) (string, error) {
	bolck, _ := pem.Decode([]byte(pub_key))
	if bolck == nil {
		return "", errors.New("E")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(bolck.Bytes)
	if err != nil {
		return "", errors.New("E:" + err.Error())
	}
	pub := pubInterface.(*rsa.PublicKey)
	sign_new, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(data))
	if err != nil {
		return "", errors.New("E:" + err.Error())
	}
	return base64.StdEncoding.EncodeToString(sign_new), err
}
func GooglePlayRsaVerifyWithSha1Base64(ori_data, sign_data, pubKey string) error {
	sign, err := base64.StdEncoding.DecodeString(sign_data)
	if err != nil {
		return err
	}
	public, _ := base64.StdEncoding.DecodeString(pubKey)
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(ori_data))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

func RsaVerifyWithPubkeyMD5Demo(originStr, rsaPubkey, signReceived string) error {
	signReceivedBase64decode, err := base64.StdEncoding.DecodeString(signReceived)
	if err != nil {
		return err
	}
	public, err := base64.StdEncoding.DecodeString(rsaPubkey)
	if err != nil {
		return err
	}

	pubInterface, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}
	hash := md5.New()
	hash.Write([]byte(originStr))
	err = rsa.VerifyPKCS1v15(pubInterface.(*rsa.PublicKey), crypto.MD5, hash.Sum(nil), signReceivedBase64decode)
	if err != nil {
		return err
	}
	return nil
}

func Md5WithRsaUsePublicKey(originStr string, pubKey string, signReceived string) error {
	signReceivedBase64decode, err := base64.StdEncoding.DecodeString(signReceived)
	if err != nil {
		return err
	}
	public, _ := base64.StdEncoding.DecodeString(pubKey)

	pubInterface, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}
	hash := md5.New()
	hash.Write([]byte(originStr))
	err = rsa.VerifyPKCS1v15(pubInterface.(*rsa.PublicKey), crypto.MD5, hash.Sum(nil), signReceivedBase64decode)
	if err != nil {
		return err
	}
	return nil
}
