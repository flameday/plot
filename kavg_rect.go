package main

func getDistanceArray(data []float64, startPos int, aimPos int) []float64 {
	arr := make([]float64, 0)
	if startPos < aimPos {
		for i := aimPos - 1; i >= startPos; i-- {
			val := data[i] - data[aimPos]
			arr = append(arr, val*val)
		}
	}
	if startPos > aimPos {
		for i := aimPos + 1; i <= startPos; i++ {
			val := data[i] - data[aimPos]
			arr = append(arr, val*val)
		}
	}
	return arr
}

func compareDistanceArray(leftArr []float64, rightArr []float64) int {
	cnt := 0
	for i := 0; i < 8; i++ {
		if i >= len(leftArr) {
			break
		}
		if i >= len(rightArr) {
			break
		}
		if leftArr[i] < rightArr[i] {
			cnt += 1
		} else {
			cnt -= 1
		}
	}
	return cnt
}
func getRectangle(data []float64, posLeft int, posMiddle int, posRight int) (int, int) {
	leftArr := getDistanceArray(data, posLeft, posMiddle)
	rightArr := getDistanceArray(data, posMiddle, posRight)
	val := compareDistanceArray(leftArr, rightArr)
	if val > 0 {
		return posLeft, posMiddle
	}
	return posMiddle, posRight
}
