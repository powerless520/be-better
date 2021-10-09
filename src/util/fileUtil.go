package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadFileInfo(filePath string) (dataBytes []byte, err error) {

	_, err = PathIfExists(filePath)
	if err != nil {
		fmt.Println("判断文件不存在", err)
		return dataBytes, err
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件失败", err)
		return bytes, err
	}

	fmt.Println("读取成功")
	return bytes, nil
}

func PathIfExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func demo(b []byte) {
	var result []string
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
}
