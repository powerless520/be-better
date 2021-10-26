package encryptUtil

import (
	"strings"
)

const Base32Digit = "0xd4uqe6fnbyms9ap187wjr5tkg32zhc"

var base32DigitMap = map[string]int{
	"0": 0, "x": 1, "d": 2, "4": 3, "u": 4, "q": 5, "e": 6, "6": 7,
	"f": 8, "n": 9, "b": 10, "y": 11, "m": 12, "s": 13, "9": 14, "a": 15,
	"p": 16, "1": 17, "8": 18, "7": 19, "w": 20, "j": 21, "r": 22, "5": 23,
	"t": 24, "k": 25, "g": 26, "3": 27, "2": 28, "z": 29, "h": 30, "c": 31,
}

func Base32Encode(data string) string {
	var (
		i, index, digit, currentByte, nextByte int
	)

	dataByte := []byte(data)
	var result = make([]string, 0, (len(dataByte)+4)*8/5)
	for i < len(dataByte) {
		currentByte = int(dataByte[i])
		if index > 3 {
			if i+1 < len(dataByte) {
				nextByte = int(dataByte[i+1])
			} else {
				nextByte = 0
			}

			digit = currentByte & (0xFF >> index)
			index = (index + 5) % 8
			digit <<= index
			digit |= nextByte >> (8 - index)
			i++
		} else {
			digit = currentByte >> (8 - (index + 5)) & 0x1F
			index = (index + 5) % 8
			if index == 0 {
				i++
			}
		}

		result = append(result, string(Base32Digit[digit]))
	}

	return strings.Join(result, "")
}

func Base32Decode(data string) []byte {
	var index = 0
	var offset = 0
	result := make([]byte, len(data)*5/8, len(data)*5/8)

	for i := 0; i < len(data); i++ {
		digit, ok := base32DigitMap[string(data[i])]
		if !ok {
			continue
		}
		current := 0
		if index <= 3 {
			index = (index + 5) % 8
			if index == 0 {
				current = int(result[offset]) | digit
				result[offset] = byte(current)
				offset++
				if offset >= len(result) {
					break
				}
			} else {
				current = int(result[offset]) | digit<<(8-index)
				result[offset] = byte(current)
			}
		} else {
			index = (index + 5) % 8
			current = int(result[offset]) | digit>>index
			result[offset] = byte(current)
			offset++
			if offset >= len(result) {
				break
			}
			current = int(result[offset]) | digit<<(8-index)
			result[offset] = byte(current)
		}
	}

	return result
}
