package main

//get by time order
func isSimilar(r1 *Rect, r2 *Rect) bool {
	height1 := r1.Height()
	height2 := r2.Height()

	if height1 > 2*height2 {
		return false
	}
	if height2 > 2*height1 {
		return false
	}
	if r1.Width() > 34 {
		return false
	}
	if r2.Width() > 34 {
		return false
	}

	if height1 >= height2 {
		if r1.top >= r2.top && r1.bottom <= r2.bottom {
			//contain
			return true
		} else if height1 > 2*height2 {
			//not 2 times
			return false
		} else if (r1.top-height1*0.618 <= r2.top) && (r1.bottom+height1*0.618 >= r2.bottom) {
			//50% is threshold
			return true
		}
		return false
	} else if height1 < height2 {
		if r2.top >= r1.top && r2.bottom <= r1.bottom {
			//contain
			return true
		} else if height2 > 2*height1 {
			//not 2 times
			return false
		} else if (r2.top-height2*0.618 <= r1.top) && (r2.bottom+height2*0.618 >= r1.bottom) {
			//50% is threshold
			return true
		}
	}
	return false
}
func mergeRect(arr []Rect) []int {
	flagArr := make([]int, len(arr))
	for i := 0; i < len(arr)-1; i++ {
		if isSimilar(&arr[i], &arr[i+1]) {
			flagArr[i] = 1
		} else {
			flagArr[i] = 0
		}
	}
	//必须连续3笔才行
	lastCnt := 0
	for i := 0; i < len(flagArr); i++ {
		if flagArr[i] == 1 {
			lastCnt += 1
		} else {
			if lastCnt > 0 && lastCnt < 2 {
				for k := 0; k < lastCnt; k++ {
					pos := i - k - 1
					flagArr[pos] = 0
				}
			}
			lastCnt = 0
		}
	}
	return flagArr
}

//get by price order
