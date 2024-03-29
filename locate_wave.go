package main

func isLow(data []float64, index int, length int) bool {
	//前后端点不计算大小值
	if index <= length/2 || index >= len(data)-length/2 {
		return false
	}

	for i := 1; i <= length/2; i++ {
		if data[index-i] < data[index] {
			return false
		}
		if data[index+i] < data[index] {
			return false
		}
	}
	return true
}
func isHigh(data []float64, index int, length int) bool {
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
	return true
}
func isRising(stock *Stock, pre int, post int) bool {
	if stock.dataMinMax[pre] < 0 && stock.dataMinMax[post] > 0 {
		return true
	}
	return false
}
func isFall(stock *Stock, pre int, post int) bool {
	if stock.dataMinMax[pre] > 0 && stock.dataMinMax[post] < 0 {
		return true
	}
	return false
}
func getSmallWave(stock *Stock, pre int, post int) {
	if isRising(stock, pre, post) {
		//先求高点
		for i := pre + 1; i < post; i++ {
			if isHigh(stock.dataHigh, i, 5) {
				stock.subDataMinMax[i] = 2
			}
		}
		//遍历高点，求后面的低点
		for i := pre + 1; i < post; i++ {
			if stock.subDataMinMax[i] == 2 {
				next := findNextIndex(stock.subDataMinMax, i+1, 2)
				if next == -1 {
					next = post
				}
				for k := i + 2; k < next; k++ {
					if isLow(stock.dataLow, k, 5) {
						//判断高低
						if stock.dataHigh[i] > stock.dataHigh[k] && stock.dataLow[i] > stock.dataLow[k] {
							stock.subDataMinMax[k] = -2
						}
					}
				}
			}
		}
	} else if isFall(stock, pre, post) {
		for i := pre + 1; i < post; i++ {
			if isLow(stock.dataLow, i, 5) {
				stock.subDataMinMax[i] = -2
			}
		}
		//遍历低点，求后面的高点
		for i := pre + 1; i < post; i++ {
			if stock.subDataMinMax[i] == -2 {
				next := findNextIndex(stock.subDataMinMax, i+1, -2)
				if next == -1 {
					next = post
				}
				for k := i + 2; k < next; k++ {
					if isHigh(stock.dataLow, k, 5) {
						//判断高低
						if stock.dataHigh[i] < stock.dataHigh[k] && stock.dataLow[i] < stock.dataLow[k] {
							stock.subDataMinMax[k] = 2
						}
					}
				}
			}
		}
	}
	//合并相邻的高点、合并相邻的低点
	lastPos := -1
	for i := pre + 1; i < post; i++ {
		if (stock.subDataMinMax[i] == 2) || (stock.subDataMinMax[i] == -2) {
			if lastPos == -1 || stock.subDataMinMax[lastPos] != stock.subDataMinMax[i] {
				lastPos = i
			} else if stock.subDataMinMax[lastPos] == stock.subDataMinMax[i] {
				if stock.subDataMinMax[lastPos] == 2 {
					if stock.dataHigh[lastPos] <= stock.dataHigh[i] {
						stock.subDataMinMax[lastPos] = 0

						lastPos = i
					} else {
						stock.subDataMinMax[i] = 0
						//keep lastPos
					}
				} else if stock.subDataMinMax[lastPos] == -2 {
					if stock.dataLow[lastPos] >= stock.dataLow[i] {
						stock.subDataMinMax[lastPos] = 0

						lastPos = i
					} else {
						stock.subDataMinMax[i] = 0
					}
				}
			}
		}
	}
}
func firstWave(stock *Stock, index int) {
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
	minMax := stock.dataMinMax

	pos, _ := findPreMinOrMaxIndex(minMax, index-1)
	if pos == -1 {
		firstWave(stock, index)
		return
	}
	preLowPos := findPreIndex(minMax, index-1, -1)
	preHighPos := findPreIndex(minMax, index-1, 1)

	highCntFromLow := 0
	lowCntFromLow := 0

	highCntFromHigh := 0
	lowCntFromHigh := 0
	tmpPreLowPos := preLowPos
	if tmpPreLowPos == -1 {
		tmpPreLowPos = 0
	}
	for i := tmpPreLowPos; i <= index; i++ {
		if stock.dataHigh[i] < stock.avg10[i] {
			lowCntFromLow += 1
		} else if stock.dataLow[i] > stock.avg10[i] {
			highCntFromLow += 1
		}
	}

	tmpPreHighPos := preHighPos
	if tmpPreHighPos == -1 {
		tmpPreHighPos = 0
	}
	for i := tmpPreHighPos; i <= index; i++ {
		if stock.dataHigh[i] < stock.avg10[i] {
			lowCntFromHigh += 1
		} else if stock.dataLow[i] > stock.avg10[i] {
			highCntFromHigh += 1
		}
	}
	//可能有一个为 -1
	//低 ---> 高
	if preLowPos < preHighPos || preLowPos == -1 {
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
	if preHighPos < preLowPos || preHighPos == -1 {
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
