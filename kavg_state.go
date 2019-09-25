package main

import "github.com/prometheus/common/log"

// 初始化
func isValidInit(ac *avgContext, data []float64) bool {
	//get pre arr
	arr, _ := getAllRect(data)
	if len(arr) <= 1 {
		return false
	}

	curValue := data[len(data)-1]

	// M
	size := len(arr)
	if (arr[size-1].top < arr[size-2].top) && (arr[size-1].bottom < arr[size-2].bottom) {
		if curValue < arr[0].bottom {
			ac.State = STATE_NEW_LOW
			ac.Action = ACTION_SELL
			ac.Sell_stop = arr[size-1]
			return true
		}
	}
	// W
	if (arr[size-1].top > arr[size-2].top) && (arr[size-1].bottom > arr[size-2].bottom) {
		if curValue > arr[size-1].top {
			ac.State = STATE_NEW_HIGH
			ac.Action = ACTION_BUY
			ac.Buy_stop = arr[size-1]
			return true
		}
	}
	return false
}

func action_High_Buy(ac *avgContext, arr []Rect, curValue float64) bool {
	size := len(arr)
	if curValue >= ac.Buy_stop.bottom {
		// 推进 Buy_stop
		if (arr[size-1].top > ac.Buy_stop.top) && (arr[size-1].bottom >= ac.Buy_stop.bottom) {
			ac.Buy_stop = arr[size-1]

			return true
		}
	} else if curValue < ac.Buy_stop.bottom {
		ac.State = STATE_NEW_HIGH + STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Sell_stop = arr[size-1]
		ac.Min_High_low = curValue

		return true
	}

	return false
}
func action_High_Sell(ac *avgContext, arr []Rect, curValue float64) bool {

	return false
}
func action_Low_Buy(ac *avgContext, arr []Rect, curValue float64) bool {
	return false
}
func action_Low_Sell(ac *avgContext, arr []Rect, curValue float64) bool {
	return false
}
func action_High_Low_Buy(ac *avgContext, arr []Rect, curValue float64) bool {
	return false
}
func action_High_Low_Sell(ac *avgContext, arr []Rect, curValue float64) bool {
	return false
}
func action_Low_High_Buy(ac *avgContext, arr []Rect, curValue float64) bool {
	return false
}
func action_Low_High_Sell(ac *avgContext, arr []Rect, curValue float64) bool {
	return false
}

func forwardStateFromBuy(ac *avgContext, data []float64) bool {

	if ac.Action != ACTION_BUY {
		log.Error("forwardStateFromBuy state error")
		return false
	}

	//get pre arr
	arr, _ := getAllRect(data)
	if len(arr) <= 1 {
		return false
	}

	curValue := data[len(data)-1]

	size := len(arr)
	if ac.State == STATE_NEW_HIGH {
		if ac.Action == ACTION_BUY {

		}
	} else if ac.State == STATE_NEW_HIGH+STATE_NEW_LOW {
		// 更新最低点
		if curValue < ac.Min_High_low {
			ac.Min_High_low = curValue
		}
		if ac.Action == ACTION_BUY {
			if curValue > ac.Sell_stop.top {
				//A
				ac.SubState = "A"
			} else if curValue > ac.Sell_stop.bottom {
				//B C
				ac.SubState = "BC"
			} else {
				//keep
			}
		} else if ac.Action == ACTION_SELL {
			if curValue <= ac.Sell_stop.bottom {
				//keep sell
			} else if curValue > ac.Sell_stop.bottom {
				ac.Action = ACTION_BUY
				//这里需要获取最低点
				ac.Buy_stop.bottom = ac.Min_High_low
				ac.Buy_stop.top = ac.Sell_stop.top
			}
		}
	}
	return true
}

func forwardStateFromSell(ac *avgContext, data []float64) bool {

	if ac.Action != ACTION_SELL {
		log.Error("forwardStateFromSell state error")
		return false
	}

	//get pre arr
	arr, _ := getAllRect(data)
	if len(arr) <= 1 {
		return false
	}

	curValue := data[len(data)-1]

	size := len(arr)

	if ac.State == STATE_NEW_LOW {
		// New Low
		if curValue > ac.Sell_stop.top {
			ac.State = STATE_NEW_LOW + STATE_NEW_HIGH
			ac.Action = ACTION_BUY
			ac.Sell_stop = arr[size-1]

			return true
		}

		// 推进 Sell_stop
		if (arr[size-1].top <= ac.Sell_stop.top) && arr[size-1].bottom < ac.Sell_stop.bottom {
			ac.Sell_stop = arr[size-1]

			return true
		}
	}
	return false
}
