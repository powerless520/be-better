package encryptUtil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"io"
)

func AesGcmEncrypt(src, key string) (string, error) {
	keyBytes, err := hex.DecodeString(key)
	plaintext := []byte(src)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	//plaintext = PKCS5Padding(plaintext, block.BlockSize())

	var n = make([]byte, 12)
	io.ReadFull(rand.Reader, n)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText := aesgcm.Seal(nil, n, plaintext, nil)
	cipherText = BytesCombine(n, cipherText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func AesCfbEncrypt(data []byte, key []byte, iv []byte) (encryptData []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	iv = iv[:aes.BlockSize]
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(data, data)
	return data, nil
}

func AesCfbDecrypt(data []byte, key []byte, iv []byte) (decryptData []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	iv = iv[:aes.BlockSize]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(data, data)
	return data, nil
}

func Base32Encrypt(data, key, iv string) (encryptData string, err error) {
	raw, err := AesCfbEncrypt([]byte(data), []byte(key), []byte(iv))
	if err != nil {
		return
	}

	sum := crc32.ChecksumIEEE([]byte(data))
	p := Protocol{Format: []string{"N4"}}
	b := p.Pack(int64(sum))
	raw = append(raw, b...)

	return Base32Encode(string(raw)), nil
}

func Base32Decrypt(data, key, iv string) (string, error) {
	originData := Base32Decode(data)
	if len(originData) < 4 {
		return "", errors.New("decrypt data is empty")
	}
	checkData := originData[len(originData)-4:]
	p := Protocol{Format: []string{"N4"}}
	checkByte := p.UnPack(checkData)
	rawData := originData[:len(originData)-4]
	decryptByte, err := AesCfbDecrypt(rawData, []byte(key), []byte(iv))
	if err != nil {
		return "", err
	}

	if checkByte[0] != int64(crc32.ChecksumIEEE(decryptByte)) {
		return "", errors.New("check sum failed")
	}

	return string(decryptByte), nil
}
