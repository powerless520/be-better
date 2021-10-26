package utils

func BubbleSort(arr []int, n int) {
	var i, j, k int
	for i = 0; i < n; i++ {
		k = 0
		for j = 0; j < n-1; j++ {
			if arr[j] > arr[j+1] {
				t := arr[j]
				arr[j] = arr[j+1]
				arr[j+1] = t
				k = 1
			}
		}
		if k == 0 {
			break
		}
	}
}
