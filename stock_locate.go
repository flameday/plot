package main

import (
	log "github.com/cihub/seelog"
	"runtime/debug"
)

func findMinIndex(data []float64, flagArr []int, posStart int, posEnd int) int {
	if posStart >= posEnd {
		return -1
	}
	minIndex := posStart
	minValue := data[minIndex]
	for i := posStart + 1; i < posEnd; i++ {
		//small
		if flagArr[i] == -1 {
			if data[i] < minValue {
				minValue = data[i]
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
func isMax(data []float64, avg []float64, index int, length int) bool {
	//前后端点不计算大小值
	if index <= length/2 || index >= len(data)-length/2 {
		return false
	}

	for i := 1; i <= length/2; i++ {
		if data[index-i] > data[index] {
			return false
		}
		if data[index+i] > data[index] {
			return false
		}
	}
	//
	limitCnt := 0
	for i := 0; i < index+3; i++ {
		if i <= 0 {
			continue
		}
		if i >= len(data) || i >= len(avg) {
			continue
		}
		if data[i] <= avg[i] {
			break
		}
		limitCnt += 1
		if limitCnt >= 3 {
			return true
		}
	}
	for i := index - 1; i >= index-3; i-- {
		if i <= 0 {
			continue
		}
		if i >= len(data) || i >= len(avg) {
			continue
		}

		if data[i] <= avg[i] {
			break
		}
		limitCnt += 1
		if limitCnt >= 3 {
			return true
		}
	}
	return false
}
func initWave(stock *Stock, index int) {
	//minMax := stock.dataMinMax
	prePos, _ := findPreMinOrMaxIndex(stock.dataMinMax, index-1)
	if prePos == -1 {
		prePos = 0
	}
	highCnt := 0
	lowCnt := 0
	for i := prePos; i <= index; i++ {
		if stock.dataHigh[i] < stock.avg10[i] {
			lowCnt += 1
		} else if stock.dataLow[i] > stock.avg10[i] {
			highCnt += 1
		}
	}
	if highCnt >= 3 {
		stock.dataMinMax[index] = 1
	} else if lowCnt >= 3 {
		stock.dataMinMax[index] = -1
	}
}
func getWave(stock *Stock, index int) {
	//if index == 21 {
	//	log.Infof("index:%d", index)
	//}

	minMax := stock.dataMinMax

	pos, _ := findPreMinOrMaxIndex(minMax, index-1)
	if pos == -1 {
		initWave(stock, index)
		return
	}
	preLowPos := findPreIndex(minMax, index-1, -1)
	preHighPos := findPreIndex(minMax, index-1, 1)

	highCntFromLow := 0
	lowCntFromLow := 0

	highCntFromHigh := 0
	lowCntFromHigh := 0
	if preLowPos != -1 {
		for i := preLowPos; i <= index; i++ {
			if stock.dataHigh[i] < stock.avg10[i] {
				lowCntFromLow += 1
			} else if stock.dataLow[i] > stock.avg10[i] {
				highCntFromLow += 1
			}
		}
	}
	if preHighPos != -1 {
		for i := preHighPos; i <= index; i++ {
			if stock.dataHigh[i] < stock.avg10[i] {
				lowCntFromHigh += 1
			} else if stock.dataLow[i] > stock.avg10[i] {
				highCntFromHigh += 1
			}
		}
	}

	//低 ---> 高
	if preLowPos < preHighPos {
		if highCntFromLow >= 3 {
			//merge
			if stock.dataHigh[preHighPos] > stock.dataHigh[index] {
				//keep old
			} else if stock.dataHigh[index] > stock.avg10[index] {
				//use new
				minMax[preHighPos] = 0
				minMax[index] = 1
			}
		}
		if lowCntFromHigh >= 3 {
			minMax[index] = -1
		}
	}
	//高 ---> 低
	if preHighPos < preLowPos {
		if lowCntFromHigh >= 3 {
			//merge
			if stock.dataLow[preLowPos] < stock.dataLow[index] {
				//keep old
			} else if stock.dataLow[index] < stock.avg10[index] {
				//use new
				minMax[preLowPos] = 0
				minMax[index] = -1
			}
		}

		if highCntFromLow >= 3 {
			minMax[index] = 1
		}
	}

}

//func getWave3(stock *Stock, index int) {
//	//if index == 21 {
//	//	log.Infof("index:%d", index)
//	//}
//
//	minMax := stock.dataMinMax
//
//	pos, _ := findPreMinOrMaxIndex(minMax, index-1)
//	if pos == -1 {
//		initWave(stock, index)
//		return
//	}
//
//	//找到1，-1的位置
//	preLowPos := findPreIndex(minMax, index-1, -1)
//	preHighPos := findPreIndex(minMax, index-1, 1)
//	// 二选一
//	if preLowPos == -1 || preHighPos == -1 {
//		if preLowPos != -1 {
//			highCnt := 0
//			lowCnt := 0
//			for i := preLowPos; i <= index; i++ {
//				if stock.dataHigh[i] < stock.avg10[i] {
//					lowCnt += 1
//				} else if stock.dataLow[i] > stock.avg10[i] {
//					highCnt += 1
//				}
//			}
//
//			if lowCnt >= 3 {
//				//merge
//				if stock.dataLow[preLowPos] < stock.dataLow[index] {
//					//keep old
//				} else if stock.dataLow[index] < stock.avg10[index] {
//					//use new
//					minMax[preLowPos] = 0
//					minMax[index] = -1
//					log.Errorf("		-1: use new %d  ---> %d", preLowPos, index)
//				}
//			}
//
//			if highCnt >= 3 {
//				minMax[index] = 1
//			}
//		}
//
//		if preHighPos != -1 {
//			highCnt := 0
//			lowCnt := 0
//			for i := preHighPos; i <= index; i++ {
//				if stock.dataHigh[i] < stock.avg10[i] {
//					lowCnt += 1
//				} else if stock.dataLow[i] > stock.avg10[i] {
//					highCnt += 1
//				}
//			}
//			if highCnt >= 3 {
//				//merge
//				if stock.dataHigh[preHighPos] > stock.dataHigh[index] {
//					//keep old
//				} else if stock.dataHigh[index] > stock.avg10[index] {
//					//use new
//					minMax[preHighPos] = 0
//					minMax[index] = 1
//					log.Errorf("		1: use new %d  ---> %d", preHighPos, index)
//				}
//			}
//			if lowCnt >= 3 {
//				minMax[index] = -1
//			}
//		}
//
//		return
//	}
//	//两个都有:先低后高
//	if preLowPos < preHighPos {
//		highCnt := 0
//		lowCnt := 0
//		for i := preLowPos; i <= index; i++ {
//			if stock.dataHigh[i] < stock.avg10[i] {
//				lowCnt += 1
//			} else if stock.dataLow[i] > stock.avg10[i] {
//				highCnt += 1
//			}
//		}
//		//合并高点
//		if highCnt >= 3 {
//			//merge
//			if stock.dataHigh[preHighPos] > stock.dataHigh[index] {
//				//keep old
//			} else if stock.dataHigh[index] > stock.avg10[index] {
//				//use new
//				minMax[preHighPos] = 0
//				minMax[index] = 1
//				log.Errorf("		1: use new %d  ---> %d", preHighPos, index)
//			}
//		}
//		if lowCnt >= 3 {
//			minMax[index] = -1
//		}
//	}
//
//	// 高低
//	if preHighPos < preLowPos {
//		highCnt := 0
//		lowCnt := 0
//		for i := preHighPos; i <= index; i++ {
//			if stock.dataHigh[i] < stock.avg10[i] {
//				lowCnt += 1
//			} else if stock.dataLow[i] > stock.avg10[i] {
//				highCnt += 1
//			}
//		}
//		//合并高点
//		if lowCnt >= 3 {
//			//merge
//			if stock.dataLow[preLowPos] < stock.dataLow[index] {
//				//keep old
//			} else if stock.dataLow[index] < stock.avg10[index] {
//				//use new
//				minMax[preLowPos] = 0
//				minMax[index] = -1
//				log.Errorf("		-1: use new %d  ---> %d", preLowPos, index)
//			}
//		}
//		if highCnt >= 3 {
//			minMax[index] = 1
//		}
//	}
//
//	if minMax[index] == 1 || minMax[index] == -1 {
//		log.Errorf("[%d] minMax:%d", index, minMax[index])
//	}
//}
func get_fractal(data []float64, index int) int {
	if index == 0 || index >= len(data)-1 {
		return 0
	}

	if data[index] > data[index-1] && data[index] > data[index+1] {
		return 1
	}
	if data[index] < data[index-1] && data[index] < data[index+1] {
		return 1
	}
	return 0
}

// 局部最小值
func isMin(data []float64, avg []float64, index int, length int) bool {
	//前后端点不计算大小值
	if index <= length/2 || index >= len(data)-length/2 {
		return false
	}

	// 判断底分型
	for i := 1; i <= length/2; i++ {
		if data[index-i] < data[index] {
			return false
		}
		if data[index+i] < data[index] {
			return false
		}
	}
	//判断 3 根K线在 10均线 之下
	limitCnt := 0
	for i := 0; i < index+3; i++ {
		if i <= 0 {
			continue
		}
		if i >= len(data) || i >= len(avg) {
			continue
		}
		if data[i] >= avg[i] {
			break
		}
		limitCnt += 1
		if limitCnt >= 3 {
			return true
		}
	}
	for i := index - 1; i >= index-3; i-- {
		if i <= 0 {
			continue
		}
		if i >= len(data) || i >= len(avg) {
			continue
		}

		if data[i] >= avg[i] {
			break
		}
		limitCnt += 1
		if limitCnt >= 3 {
			return true
		}
	}
	return false
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
func locateMax(dataClose []float64, posArr []int, length int) {
	defer func() {
		if err := recover(); err != nil {
			log.Infof("error")
			defer log.Flush()
			debug.PrintStack()
		}
	}()

	//for max
	for i := 0; i < len(posArr); i++ {
		//log.Infof("[%d] dataClose:%f", i, dataClose[i])
		if posArr[i] == MAX_VALUE_FLAG {
			//max
			maxValue := dataClose[i]
			maxPos := i
			//right to left
			for j := length / 2; j >= -1*length/2; j-- {
				pos := i + j
				if pos >= len(posArr) || pos < 0 {
					continue
				}
				//log.Infof("maxValue:%.2f len:%d pos:%d dataClose[pos]:%f", maxValue, len(dataClose), pos, dataClose[pos])
				if maxValue < dataClose[pos] {

					posArr[maxPos] = 0
					posArr[pos] = MAX_VALUE_FLAG

					maxValue = dataClose[pos]
					maxPos = pos
				} else if maxValue == dataClose[pos] {
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
func locateMin(dataClose []float64, posArr []int, length int) {
	//for max
	if len(posArr) != len(dataClose) {
		log.Errorf("not equal")
	}

	for i := 0; i < len(posArr); i++ {
		//log.Infof("[%d] dataClose:%f", i, dataClose[i])
		if posArr[i] == MIN_VALUE_FLAG {
			//max
			minValue := dataClose[i]
			minPos := i
			//right to left
			for j := length / 2; j >= -1*length/2; j-- {
				pos := i + j
				if pos >= len(posArr) || pos < 0 {
					continue
				}
				if minValue > dataClose[pos] {
					//log.Infof("[%d] pos dataClose:%f", pos, dataClose[pos])

					posArr[minPos] = 0
					posArr[pos] = MIN_VALUE_FLAG

					minValue = dataClose[pos]
					minPos = pos
				} else if minValue == dataClose[pos] {
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
