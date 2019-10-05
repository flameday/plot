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
	for i := posStart; i >= 0 && i < len(flagArr); i-- {
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

func findPreMinOrMaxIndex(flagArr []int, posStart int) (int, int) {
	preMin := findPreIndex(flagArr, posStart, -1)
	preMax := findPreIndex(flagArr, posStart, 1)
	if preMin == -1 && preMax == -1 {
		return -1, -1
	}
	if preMin == -1 {
		return preMax, 1
	}
	if preMax == -1 {
		return preMin, -1
	}
	if preMin > preMax {
		return preMin, -1
	}
	return preMax, 1
}
