package test

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func TestDemo(t *testing.T) {

	fmt.Println(1<<2)

}

// GetTableNameFromOrderId 根据订单号计算order表名
func GetTableNameFromOrderId(order_id string) string {
	unixtime_str_int64 := time.Now().Unix()
	var err error
	if err = VerifyValidOrderId(order_id); err == nil {
		unixtime_str := order_id[len(order_id)-10 : len(order_id)]
		unixtime_str_int64, err = strconv.ParseInt(unixtime_str, 10, 64)
		if err != nil {
			unixtime_str_int64 = time.Now().Unix()
		}
	}
	unixtime := time.Unix(unixtime_str_int64, 0)
	year, week := unixtime.ISOWeek()
	year_week := strconv.Itoa(year*100 + week)
	return "t_order_" + year_week
}

func VerifyValidOrderId(order_id string) (err error) {
	if order_id == "" {
		return errors.New("订单号为空")
	}
	if len(order_id) < 10 {
		return errors.New("订单号错误")
	}
	return nil
}

/*公钥解密*/
func pubKeyDecrypt(pub *rsa.PublicKey, data []byte) []byte {
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	//d := leftPad(m.Bytes(), k)
	d := make([]byte,2)
	if d[0] != 0 {
		return nil
	}
	if d[1] != 0 && d[1] != 1 {
		return nil
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil
	}
	return d[i:]
}

// RsaDecrypt 使用公钥解密
func RsaDecrypt(data string, publicKey string, isFomartKey bool) (decryptData string, err error) {
	if len(data) <= 0 || len(publicKey) <= 0 {
		return
	}

	if isFomartKey {
		//publicKey = formatKey(publicKey)
	}

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub := pubInterface.(*rsa.PublicKey)

	retDataBytes, err := base64.StdEncoding.DecodeString(data)
	//decLen = 128 分段解密
	decLen := pub.N.BitLen() / 8
	decchunks := split(retDataBytes, decLen)
	decbuffer := bytes.NewBufferString("")
	for _, decchunk := range decchunks {
		decbytes := pubKeyDecrypt(pub, decchunk)
		decbuffer.Write(decbytes)
	}

	decryptData = string(decbuffer.Bytes())
	//git上封装的解密方法 一句话搞定
	//decrypt, err := gorsa.PublicDecrypt(responseData["data"].(string), formatPubkey)

	return
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
