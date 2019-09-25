package main

func isValidInit(data []float64, curPos int) (bool, string, int) {
	//get pre arr
	arr, _ := getAllRect(data)
	if len(arr) <= 1 {
		return false, "", -1
	}
	size := len(arr)

	if (arr[size-1].top < arr[size-2].top) && (arr[size-1].bottom < arr[size-2].bottom) {
		if data[curPos] < arr[0].bottom {
			return true, STATE_NEW_LOW, ACTION_SELL
		}
	}
	if (arr[size-1].top > arr[size-2].top) && (arr[size-1].bottom > arr[size-2].bottom) {
		if data[curPos] > arr[size-1].top {
			return true, STATE_NEW_HEIGHT, ACTION_BUY
		}
	}
	return false, "", -1
}
