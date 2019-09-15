package main

import "github.com/prometheus/common/log"

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

// 根据前面的 length 个值，获取平均值
func get_pre_avg(value_list []float64, index int, length int) float64 {
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

//根据前后 length 个值，获取平均值
func get_middle_avg(value_list []float64, index int, length int) float64 {
	var total float64
	var count int
	total = 0.0
	count = 0
	for i := 0; i < length/2; i++ {
		pos := index - i
		if pos < 0 {
			break
		}
		total += value_list[pos]
		count += 1
	}
	for i := 0; i < length/2; i++ {
		pos := index + i
		if pos < 0 || pos >= len(value_list) {
			break
		}
		total += value_list[pos]
		count += 1
	}
	return total / float64(count)
}

// 向前向后length找最大值
func locateMax(data []float64, posArr []int, length int) {
	//for max
	for i := 0; i < len(posArr); i++ {
		log.Infof("[%d] data:%f", i, data[i])
		if posArr[i] == MAX_VALUE_FLAG {
			//max
			maxValue := data[i]
			maxPos := i
			//right to left
			for j := length / 2; j >= -1*length/2; j-- {
				pos := i + j
				if pos >= len(posArr) || pos < 0 {
					continue
				}
				if maxValue < data[pos] {
					log.Infof("[%d] pos data:%f", pos, data[pos])

					posArr[maxPos] = 0
					posArr[pos] = MAX_VALUE_FLAG

					maxValue = data[pos]
					maxPos = pos
				} else if maxValue == data[pos] {
					if pos > maxPos {
						//取前面的
						posArr[maxPos] = MAX_VALUE_FLAG
						posArr[pos] = 0
					} else {
						posArr[maxPos] = 0

						posArr[pos] = MAX_VALUE_FLAG
						maxPos = pos
					}
				}
			}
		}
	}
}

// 向前向后length找最小值
func locateMin(data []float64, posArr []int, length int) {
	//for max
	for i := 0; i < len(posArr); i++ {
		log.Infof("[%d] data:%f", i, data[i])
		if posArr[i] == MIN_VALUE_FLAG {
			//max
			minValue := data[i]
			minPos := i
			//right to left
			for j := length / 2; j >= -1*length/2; j-- {
				pos := i + j
				if pos >= len(posArr) || pos < 0 {
					continue
				}
				if minValue > data[pos] {
					log.Infof("[%d] pos data:%f", pos, data[pos])

					posArr[minPos] = 0
					posArr[pos] = MIN_VALUE_FLAG

					minValue = data[pos]
					minPos = pos
				} else if minValue == data[pos] {
					if pos > minPos {
						//取前面的
						posArr[minPos] = MIN_VALUE_FLAG
						posArr[pos] = 0
					} else {
						posArr[minPos] = 0

						posArr[pos] = MIN_VALUE_FLAG
						minPos = pos
					}
				}
			}
		}
	}
}
