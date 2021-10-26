package utils

// Intersect 交集
func Intersect(slice1 []string, slice2 []string) []string { // 取两个切片的交集
	m := make(map[string]int)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}
	return n
}

// Difference 差集
func Difference(slice1, slice2 []string) []string { //取要校验的和已经校验过的差集，找出需要校验的切片IP（找出slice1中  slice2中没有的）
	m := make(map[string]int)
	n := make([]string, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range slice1 {
		if m[value] == 0 {
			n = append(n, value)
		}
	}

	for _, v := range slice2 {
		if m[v] == 0 {
			n = append(n, v)
		}
	}
	return n
}
