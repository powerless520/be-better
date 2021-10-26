package strUtil

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func Utf8ToGbk(word string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(word)), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "Utf8ToGbkError:" + e.Error() + "|" + word
	}
	return string(d)
}
