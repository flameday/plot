package main

func findNextIndex(arr []int, posStart int, value int) int {
	for i := posStart; i < len(arr); i++ {
		if arr[i] == value {
			return i
		}
	}
	return -1
}
func findMinIndex(arr []int, posStart int, posEnd int) int {
	if posStart >= posEnd {
		return -1
	}
	minIndex := posStart
	minValue := stock.data[minIndex]
	for i := posStart + 1; i < posEnd; i++ {
		//small
		if stock.dataMinMax[i] == -1 {
			if stock.data[i] < minValue {
				minValue = stock.data[i]
				minIndex = i
			}
		}
	}
	return minIndex
}

// 中间点
//func filterMinDot(posMax []int, valueMax []float64, posMin []int, valueMin []float64) ([]int, []float64) {
//
//	posA := findNextIndex(stock.dataMinMax, -1, 1)
//	posB := findNextIndex(stock.dataMinMax, posA + 1, 1)
//	if posA != -1 && posB != -1 {
//		minIndex := findMinIndex(stock.dataMinMax, posA + 1, posB)
//	}
//}

// 局部最大值
func isMax(value_list []float64, index int, length int) bool {
	//前后端点不计算大小值
	if index <= length/2 || index >= len(value_list)-length/2 {
		return false
	}

	for i := 1; i <= length/2; i++ {
		if value_list[index-i] > value_list[index] {
			return false
		}
		if value_list[index+i] > value_list[index] {
			return false
		}
	}
	return true
}

// 局部最小值
func isMin(value_list []float64, index int, length int) bool {
	//前后端点不计算大小值
	if index <= length/2 || index >= len(value_list)-length/2 {
		return false
	}

	for i := 1; i <= length/2; i++ {
		if value_list[index-i] < value_list[index] {
			return false
		}
		if value_list[index+i] < value_list[index] {
			return false
		}
	}
	return true
}

// 获取平均值
func get_avg(value_list []float64, index int, length int) float64 {
	var total float64
	var count int
	total = 0.0
	count = 0
	for i := 0; i < length; i++ {
		pos := index - i
		if pos < 0 {
			break
		}
		total += value_list[pos]
		count += 1
	}
	return total / float64(count)
}
