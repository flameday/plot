package main

func mergeRect(arr []Rect) []int {
	flagArr := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		flagArr[i] = 0
	}
	return flagArr
}
