package randUtil

import (
	"errors"
	"math/rand"
)

func RandInt(min, max int) int {
	if min >= max {
		return max
	}
	return rand.Intn(max-min) + min
}

func RandomArray(array []int) ([]int, error) {
	if len(array) <= 0 {
		return nil, errors.New("the length of the parameter strings should not be less than 0")
	}

	for i := len(array) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		array[i], array[num] = array[num], array[i]
	}

	return array, nil
}

func IfProbability(probability float64) bool  {
	m := int(probability * 1000000)
	return RandInt(0, 1000000) <= m
}