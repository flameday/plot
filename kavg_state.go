package main

import (
	log "github.com/cihub/seelog"
)

// 初始化
func isValidInit(ac *avgContext, data []float64) (ret bool, revert bool, modify bool) {
	//get pre arr
	arr, _ := getAllRect(data)
	if len(arr) <= 1 {
		return false, false, false
	}

	curValue := data[len(data)-1]

	// M
	size := len(arr)
	if (arr[size-1].top < arr[size-2].top) && (arr[size-1].bottom < arr[size-2].bottom) {
		if curValue < arr[0].bottom {
			ac.State = STATE_NEW_LOW
			ac.Action = ACTION_SELL
			ac.Sell_stop = arr[size-1]

			return true, true, true
		}
	}
	// W
	if (arr[size-1].top > arr[size-2].top) && (arr[size-1].bottom > arr[size-2].bottom) {
		if curValue > arr[size-1].top {
			ac.State = STATE_NEW_HIGH
			ac.Action = ACTION_BUY
			ac.Buy_stop = arr[size-1]

			return true, true, true
		}
	}

	return true, false, false
}

func action_High_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	size := len(arr)
	// 正常
	if curValue >= ac.Buy_stop.bottom {
		// 推进 Buy_stop
		if (arr[size-1].top > ac.Buy_stop.top) && (arr[size-1].bottom >= ac.Buy_stop.bottom) {
			ac.Buy_stop = arr[size-1]

			return true, false, true
		}
	} else if curValue < ac.Buy_stop.bottom {
		//新低
		ac.State = STATE_NEW_HIGH__NEW_LOW_0
		ac.Action = ACTION_SELL
		ac.Sell_stop = ac.Buy_stop
		// 记录最小值
		ac.High_Low_Min = curValue

		return true, true, true
	}

	return true, false, false
}

func action_Low_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	size := len(arr)

	// 正常
	if curValue <= ac.Sell_stop.top {
		// 推进 Sell_stop
		if (arr[size-1].top <= ac.Sell_stop.top) && arr[size-1].bottom < ac.Sell_stop.bottom {
			ac.Sell_stop = arr[size-1]

			return true, false, true
		}
	}
	// 新低
	if curValue > ac.Sell_stop.top {
		ac.State = STATE_NEW_LOW__NEW_HIGH_0
		ac.Action = ACTION_BUY
		ac.Buy_stop = ac.Sell_stop
		ac.Low_High_Max = curValue

		return true, true, true
	}

	return true, false, false
}
func action_Low_High_0_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	if ac.Low_High_Max < curValue {
		ac.Low_High_Max = curValue
	}

	if curValue < ac.Buy_stop.bottom {
		//"A"
		ac.State = STATE_NEW_LOW
		// 特殊情况，一步跳到位
		if ac.Action == ACTION_BUY {
			ac.Action = ACTION_SELL
			ac.Buy_stop.bottom = ac.Sell_stop.bottom
			ac.Buy_stop.top = ac.Low_High_Max
			return true, true, true
		}
		return true, false, false
	} else if curValue < ac.Buy_stop.top {
		//
		ac.State = STATE_NEW_LOW__NEW_HIGH_1
		ac.Action = ACTION_SELL
		ac.Sell_stop.top = ac.Low_High_Max
		ac.Sell_stop.bottom = ac.Buy_stop.bottom

		return true, true, true
	}
	// 添加一个维持状态, 推进止盈
	if arr[len(arr)-1].top > ac.Buy_stop.top {
		ac.Buy_stop = arr[len(arr)-1]
	}

	return true, false, false
}

func action_Low_High_1_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	if curValue < ac.Sell_stop.bottom {
		//"A"
		ac.State = STATE_NEW_LOW

		return true, false, true
	} else if curValue > ac.Sell_stop.top {
		//"BC"
		//in fact, no "C"
		ac.State = STATE_NEW_HIGH
		ac.Action = ACTION_BUY
		ac.Buy_stop = ac.Sell_stop

		return true, true, true
	}

	return true, false, false
}

func action_High_Low_0_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// 更新最小值
	if curValue < ac.High_Low_Min {
		ac.High_Low_Min = curValue
	}

	if curValue > ac.Sell_stop.top {
		// "A"
		// 这个得排在前面，优先级更高
		ac.State = STATE_NEW_HIGH
		// 特殊情况，一步跳到位
		if ac.Action == ACTION_SELL {
			ac.Action = ACTION_BUY
			ac.Buy_stop.top = ac.Sell_stop.top
			ac.Buy_stop.bottom = ac.High_Low_Min

			return true, true, true
		}

		return true, false, true
	} else if curValue > ac.Sell_stop.bottom {
		//"BC"
		// in fact, no "C"
		ac.State = STATE_NEW_HIGH__NEW_LOW_1
		ac.Action = ACTION_BUY
		ac.Buy_stop.top = ac.Sell_stop.top
		ac.Buy_stop.bottom = ac.High_Low_Min

		return true, true, true
	}

	// 添加一个维持状态, 推进止盈
	if arr[len(arr)-1].bottom < ac.Sell_stop.bottom {
		ac.Sell_stop = arr[len(arr)-1]
	}

	return true, false, false
}

func action_High_Low_1_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// 更新最小值
	if curValue < ac.High_Low_Min {
		ac.High_Low_Min = curValue
	}

	if curValue > ac.Sell_stop.top {
		// "A"
		// 这个得排在前面，优先级更高
		ac.State = STATE_NEW_HIGH

		return true, false, false
	} else if curValue < ac.Sell_stop.bottom {
		//"BC"
		// in fact, no "C"
		ac.State = STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Sell_stop = ac.Buy_stop

		return true, true, true
	}

	return true, false, false
}

func forwardState(ac *avgContext, data []float64) (ret bool, revert bool, modify bool) {
	//get pre arr
	arr, _ := getAllRect(data)
	if len(arr) <= 1 {
		return false, false, false
	}

	ret, revert, modify = true, false, false
	curValue := data[len(data)-1]

	if ac.State == STATE_NEW_HIGH && ac.Action == ACTION_BUY {
		ret, revert, modify = action_High_Buy(ac, arr, curValue)

	} else if ac.State == STATE_NEW_HIGH && ac.Action == ACTION_SELL {
		log.Errorf("STATE_NEW_HIGH + sell impossible...")

	} else if ac.State == STATE_NEW_LOW && ac.Action == ACTION_BUY {
		log.Errorf("STATE_NEW_LOW + buy impossible...")

	} else if ac.State == STATE_NEW_LOW && ac.Action == ACTION_SELL {
		ret, revert, modify = action_Low_Sell(ac, arr, curValue)

	} else if ac.State == STATE_NEW_HIGH__NEW_LOW_0 && ac.Action == ACTION_BUY {
		log.Errorf("STATE_NEW_HIGH__NEW_LOW_0 + buy, impossible")

	} else if ac.State == STATE_NEW_HIGH__NEW_LOW_0 && ac.Action == ACTION_SELL {
		ret, revert, modify = action_High_Low_0_Sell(ac, arr, curValue)

	} else if ac.State == STATE_NEW_HIGH__NEW_LOW_1 && ac.Action == ACTION_BUY {
		ret, revert, modify = action_High_Low_1_Buy(ac, arr, curValue)

	} else if ac.State == STATE_NEW_HIGH__NEW_LOW_1 && ac.Action == ACTION_SELL {
		log.Errorf("STATE_NEW_HIGH__NEW_LOW_1 + sell, impossible")

	} else if ac.State == STATE_NEW_LOW__NEW_HIGH_0 && ac.Action == ACTION_BUY {
		ret, revert, modify = action_Low_High_0_Buy(ac, arr, curValue)

	} else if ac.State == STATE_NEW_LOW__NEW_HIGH_0 && ac.Action == ACTION_SELL {
		log.Errorf("STATE_NEW_HIGH__NEW_LOW_1 + sell, impossible")

	} else if ac.State == STATE_NEW_LOW__NEW_HIGH_1 && ac.Action == ACTION_BUY {
		log.Errorf("STATE_NEW_LOW__NEW_HIGH_1 + buy, impossible")

	} else if ac.State == STATE_NEW_LOW__NEW_HIGH_1 && ac.Action == ACTION_SELL {
		action_Low_High_1_Sell(ac, arr, curValue)
	}
	return ret, revert, modify
}
