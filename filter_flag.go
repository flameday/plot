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

func findPreIndex(flagArr []int, posStart int, flagValue int) int {
	for i := posStart; i >= 0; i-- {
		if flagArr[i] == flagValue {
			return i
		}
	}
	return -1
}

func findPostMinOrMaxIndex(flagArr []int, posStart int) int {
	postMin := findNextIndex(flagArr, posStart, -1)
	postMax := findNextIndex(flagArr, posStart, 1)
	if postMin == -1 && postMax == -1 {
		return -1
	}
	if postMin == -1 {
		return postMax
	}
	if postMax == -1 {
		return postMin
	}
	if postMin > postMax {
		return postMax
	}
	return postMin
}

func findPreMinOrMaxIndex(flagArr []int, posStart int) int {
	preMin := findPreIndex(flagArr, posStart, -1)
	preMax := findPreIndex(flagArr, posStart, 1)
	if preMin == -1 && preMax == -1 {
		return -1
	}
	if preMin == -1 {
		return preMax
	}
	if preMax == -1 {
		return preMin
	}
	if preMin > preMax {
		return preMin
	}
	return preMax
}

func filter_min_max(data []float64, minMaxArr []int) {
	last_minPos := -1
	last_maxPos := -1
	for i := 0; i < len(minMaxArr); i++ {
		if minMaxArr[i] == MIN_VALUE_FLAG {
			if last_minPos == -1 {
				last_minPos = i
			} else if last_minPos > last_maxPos {
				//判断哪个大
				if data[last_minPos] < data[i] {
					minMaxArr[i] = 0
				} else {
					minMaxArr[last_minPos] = 0
					last_minPos = i
				}
			} else {
				last_minPos = i
			}
		}
		if minMaxArr[i] == MAX_VALUE_FLAG {
			if last_maxPos == -1 {
				last_maxPos = i
			} else if last_maxPos > last_minPos {
				if data[last_maxPos] > data[i] {
					minMaxArr[i] = 0
				} else {
					minMaxArr[last_maxPos] = 0
					last_maxPos = i
				}
			} else {
				last_maxPos = i
			}
		}
	}
}

//// 2个低点之间，只能有1个高点
//func filter_max(data []float64, flagArr []int) {
//	for i := 0; i < len(flagArr); i++ {
//		if flagArr[i] == MIN_VALUE_FLAG {
//			//找到下一个小点
//			posMin := findNextIndex(flagArr, i+1, MIN_VALUE_FLAG)
//			if posMin == -1 {
//				continue
//			}
//			//找到下一个大点
//			posMax := findNextIndex(flagArr, i+1, MAX_VALUE_FLAG)
//			if posMax == -1 || posMax >= posMin {
//				continue
//			}
//			for j := posMax + 1; j < posMin; j++ {
//				if flagArr[j] == MAX_VALUE_FLAG {
//					//判断大小
//					if data[j] > data[posMax] {
//						flagArr[posMax] = 0
//						flagArr[j] = MAX_VALUE_FLAG
//
//						posMax = j
//					} else if data[j] <= data[posMax] {
//						flagArr[j] = 0
//					}
//				}
//			}
//		}
//
//	}
//
//}
//
//// 2个高点之间，只能有1个低点
//func filter_min(data []float64, flagArr []int) {
//	for i := 0; i < len(flagArr); i++ {
//		if flagArr[i] == MAX_VALUE_FLAG {
//			//找到下一个大点
//			posMin := findNextIndex(flagArr, i+1, MAX_VALUE_FLAG)
//			if posMin == -1 {
//				continue
//			}
//			//找到下一个小点
//			posMax := findNextIndex(flagArr, i+1, MIN_VALUE_FLAG)
//			if posMax == -1 || posMax >= posMin {
//				continue
//			}
//			for j := posMax + 1; j < posMin; j++ {
//				if flagArr[j] == MIN_VALUE_FLAG {
//					//判断大小
//					if data[j] < data[posMax] {
//						flagArr[posMax] = 0
//						flagArr[j] = MIN_VALUE_FLAG
//
//						posMax = j
//					} else if data[j] >= data[posMax] {
//						flagArr[j] = 0
//					}
//				}
//			}
//		}
//
//	}
//
//}
