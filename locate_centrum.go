package main

import (
	log "github.com/cihub/seelog"
)

//first time: 1/3 height
//second time: 1/3
//get by time order
func isSimilar(r1 *Rect, r2 *Rect) bool {
	height1 := r1.Height()
	height2 := r2.Height()
	log.Infof("height1: %.3f, height2:%.3f", height1, height2)

	if height1 > 3*height2 {
		return false
	}
	if height2 > 3*height1 {
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
		} else if (r1.top-height1*0.5 <= r2.top) && (r1.bottom+height1*0.5 >= r2.bottom) {
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
		} else if (r2.top-height2*0.5 <= r1.top) && (r2.bottom+height2*0.5 >= r1.bottom) {
			//50% is threshold
			return true
		}
	}
	return false
}

// merge similar rect
func mergeRect(arr []Rect) []int {
	flagArr := make([]int, len(arr))
	for i := 0; i < len(arr)-1; i++ {
		if i == 4 {
			log.Infof("debug")
		}
		if isSimilar(&arr[i], &arr[i+1]) {
			if i-1 >= 0 && (flagArr[i-1] == 1) {
				if isSimilar(&arr[i-1], &arr[i+1]) {
					flagArr[i] = 1
				}
				log.Infof("%d vs %d: %d", i-1, i+1, flagArr[i])
			} else {
				flagArr[i] = 1
			}
		}
		log.Infof("%d vs %d: %d", i, i+1, flagArr[i])
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

// check is min or max
func isMax(arr []Rect, pos int) bool {
	for i := pos - 1; i >= 0 && i >= pos-2; i-- {
		if arr[pos].top < arr[i].top {
			return false
		}
	}
	for i := pos + 1; i < len(arr) && i <= pos+2; i++ {
		if arr[pos].top < arr[i].top {
			return false
		}
	}
	return true
}
func isMin(arr []Rect, pos int) bool {
	for i := pos - 1; i >= 0 && i >= pos-2; i-- {
		if arr[pos].bottom > arr[i].bottom {
			return false
		}
	}
	for i := pos + 1; i < len(arr) && i <= pos+2; i++ {
		if arr[pos].bottom > arr[i].bottom {
			return false
		}
	}
	return true
}

//get by price order
func findExtremePoint(arr []Rect) (min int, max int) {
	min = -1
	max = -1
	//max
	for i := len(arr) - 1; i >= 0; i-- {
		if isMax(arr, i) {
			if arr[i].leftFlag == -1 {
				max = int(arr[i].right)
			} else {
				max = int(arr[i].left)
			}
		}
	}
	//min
	for i := len(arr) - 1; i >= 0; i-- {
		if isMax(arr, i) {
			if arr[i].leftFlag == -1 {
				min = int(arr[i].left)
			} else {
				min = int(arr[i].right)
			}
		}
	}
	return
}
