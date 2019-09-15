package main

//找到下一个点
func findNextIndex(flagArr []int, posStart int, flagValue int) int {
	for i := posStart; i < len(flagArr); i++ {
		if flagArr[i] == flagValue {
			return i
		}
	}
	return -1
}

// 2个低点之间，只能有1个高点
func filter_max(data []float64, flagArr []int) {
	for i := 0; i < len(flagArr); i++ {
		if flagArr[i] == MIN_VALUE_FLAG {
			//找到下一个小点
			posMin := findNextIndex(flagArr, i+1, MIN_VALUE_FLAG)
			if posMin == -1 {
				continue
			}
			//找到下一个大点
			posMax := findNextIndex(flagArr, i+1, MAX_VALUE_FLAG)
			if posMax == -1 || posMax >= posMin {
				continue
			}
			for j := posMax + 1; j < posMin; j++ {
				if flagArr[j] == MAX_VALUE_FLAG {
					//判断大小
					if data[j] > data[posMax] {
						flagArr[posMax] = 0
						flagArr[j] = MAX_VALUE_FLAG

						posMax = j
					} else if data[j] <= data[posMax] {
						flagArr[j] = 0
					}
				}
			}
		}

	}

}

// 2个高点之间，只能有1个低点
func filter_min(data []float64, flagArr []int) {
	for i := 0; i < len(flagArr); i++ {
		if flagArr[i] == MAX_VALUE_FLAG {
			//找到下一个大点
			posMin := findNextIndex(flagArr, i+1, MAX_VALUE_FLAG)
			if posMin == -1 {
				continue
			}
			//找到下一个小点
			posMax := findNextIndex(flagArr, i+1, MIN_VALUE_FLAG)
			if posMax == -1 || posMax >= posMin {
				continue
			}
			for j := posMax + 1; j < posMin; j++ {
				if flagArr[j] == MIN_VALUE_FLAG {
					//判断大小
					if data[j] < data[posMax] {
						flagArr[posMax] = 0
						flagArr[j] = MIN_VALUE_FLAG

						posMax = j
					} else if data[j] >= data[posMax] {
						flagArr[j] = 0
					}
				}
			}
		}

	}

}
