package main

import (
	"fmt"
	log "github.com/cihub/seelog"
)

type Rect struct {
	left      float64
	top       float64
	right     float64
	bottom    float64
	leftFlag  int
	rightFlag int
}

type avgContext struct {
	State          string
	Action         string
	Sell_stop      Rect
	Buy_stop       Rect
	High_Low_Min   float64
	Low_High_Max   float64
	Sell_Min_Value float64
	Buy_Max_Value  float64
	profit         float64
	buy            float64
	sell           float64
}

func (ac *avgContext) Show() string {
	s := "State : " + ac.State + " "
	s += "Action: " + ac.Action + " "
	if ac.Action == ACTION_BUY {
		s += fmt.Sprintf(" (%d, %.2f)->(%d, %.2f)",
			int(ac.Buy_stop.left),
			ac.Buy_stop.top,
			int(ac.Buy_stop.right),
			ac.Buy_stop.bottom)
	} else if ac.Action == ACTION_SELL {
		s += fmt.Sprintf(" (%d, %.2f)->(%d, %.2f)",
			int(ac.Sell_stop.left),
			ac.Sell_stop.top,
			int(ac.Sell_stop.right),
			ac.Sell_stop.bottom)
	} else {
		s += " Invalid Stop"
	}
	return s
}

// 初始化
func isValidInit(ac *avgContext, stock *Stock) (ret bool, revert bool, modify bool) {
	//get pre arr
	curIndex := len(stock.dataClose) - 1

	//简单根据高低点来尽快买入、卖出: 注意MinMax的计算方法
	preMin := findPreIndex(stock.dataMinMax, curIndex-1, -1)
	preMax := findPreIndex(stock.dataMinMax, curIndex-1, 1)
	if preMin == -1 && preMax == -1 {
		return false, false, false
	}
	if preMax == -1 {
		ac.State = STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Sell_stop = Rect{
			left:   float64(0),
			top:    stock.dataClose[0],
			right:  float64(curIndex),
			bottom: stock.dataClose[curIndex],
		}

		ac.sell = stock.dataClose[curIndex]
		ac.Buy_Max_Value = -1

		return true, true, true
	} else if preMin == -1 {
		ac.State = STATE_NEW_HIGH
		ac.Action = ACTION_BUY
		ac.Buy_stop = Rect{
			left:   float64(0),
			top:    stock.dataClose[0],
			right:  float64(curIndex),
			bottom: stock.dataClose[curIndex],
		}

		ac.buy = stock.dataClose[curIndex]
		ac.Sell_Min_Value = -1

		return true, true, true
	}

	return true, false, false
}

func action_High_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	//获取最高点
	if ac.Buy_Max_Value < curValue || ac.Buy_Max_Value == -1 {
		ac.Buy_Max_Value = curValue
	}
	// 正常
	if curValue > ac.Buy_Max_Value {
		// 推进 Buy_stop
		// 推进止损:找到前低，然后进行处理

		for size := len(arr) - 1; size >= 0; size-- {
			if arr[size].leftFlag == 1 && arr[size].rightFlag == -1 {
				if arr[size].bottom >= ac.Buy_stop.bottom && curValue > arr[size].top {
					ac.Buy_stop = arr[size]

					return true, false, true
				}
				//就看最近的回调
				break
			}
		}

	} else if curValue < ac.Buy_stop.bottom {
		//新低
		ac.State = STATE_NEW_HIGH__NEW_LOW_0
		ac.Action = ACTION_SELL
		ac.Sell_stop = ac.Buy_stop
		if ac.Buy_Max_Value != -1 {
			ac.Sell_stop.top = ac.Buy_Max_Value
			ac.Buy_Max_Value = -1
		}
		// 记录最小值
		ac.High_Low_Min = curValue
		ac.Buy_Max_Value = -1

		ac.sell = curValue
		ac.profit += curValue - ac.buy

		return true, true, true
	}

	return true, false, false
}

func action_Low_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	//获取最低点
	if ac.Sell_Min_Value > curValue || ac.Sell_Min_Value == -1 {
		ac.Sell_Min_Value = curValue
	}
	// 正常
	if curValue < ac.Sell_Min_Value {
		// 推进 Sell_stop
		for size := len(arr) - 1; size >= 0; size-- {
			if arr[size].leftFlag == 1 && arr[size].rightFlag == -1 {
				if arr[size].top <= ac.Sell_stop.top && curValue > arr[size].top {
					ac.Sell_stop = arr[size]

					return true, false, true
				}
				//就看最近的回调
				break
			}
		}
	}
	// 新低
	if curValue > ac.Sell_stop.top {
		ac.State = STATE_NEW_LOW__NEW_HIGH_0
		ac.Action = ACTION_BUY
		ac.Buy_stop = ac.Sell_stop
		if ac.Sell_Min_Value != -1 {
			ac.Buy_stop.bottom = ac.Sell_Min_Value
		}
		ac.Low_High_Max = curValue
		ac.Sell_Min_Value = -1

		ac.buy = curValue
		ac.profit += ac.sell - curValue

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

			ac.sell = curValue
			ac.profit += curValue - ac.buy
			ac.Buy_Max_Value = -1

			return true, true, true
		}
		return true, false, false
	} else if curValue < ac.Buy_stop.top {
		//"BC"
		ac.State = STATE_NEW_LOW__NEW_HIGH_1
		ac.Action = ACTION_SELL
		ac.Sell_stop.top = ac.Low_High_Max
		ac.Sell_stop.bottom = ac.Buy_stop.bottom

		ac.sell = curValue
		ac.profit += curValue - ac.buy
		ac.Buy_Max_Value = -1

		return true, true, true
	}
	// 添加一个维持状态, 推进止盈
	for size := len(arr) - 1; size >= 0; size-- {
		if arr[size].leftFlag == 1 && arr[size].rightFlag == -1 {
			if arr[size].bottom >= ac.Buy_stop.bottom && curValue > arr[size].top {
				ac.Buy_stop = arr[size]

				return true, false, true
			}
			//就看最近的回调
			break
		}
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

		ac.buy = curValue
		ac.profit += ac.sell - curValue
		ac.Sell_Min_Value = -1

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

			ac.buy = curValue
			ac.profit += ac.sell - curValue
			ac.Sell_Min_Value = -1

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

		ac.buy = curValue
		ac.profit += ac.sell - curValue
		ac.Sell_Min_Value = -1

		return true, true, true
	}

	// 添加一个维持状态, 推进止盈
	for size := len(arr) - 1; size >= 0; size-- {
		if arr[size].leftFlag == 1 && arr[size].rightFlag == -1 {
			if arr[size].top <= ac.Sell_stop.top && curValue > arr[size].top {
				ac.Sell_stop = arr[size]

				return true, false, true
			}
			//就看最近的回调
			break
		}
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

		ac.sell = curValue
		ac.profit += curValue - ac.buy
		ac.Buy_Max_Value = -1

		return true, true, true
	}

	return true, false, false
}

func forwardState(ac *avgContext, stock *Stock) (ret bool, revert bool, modify bool, arr []Rect) {
	//get pre arr
	arr, _ = getAllRect(stock)
	if len(arr) <= 1 {
		return false, false, false, arr
	}

	ret, revert, modify = true, false, false
	curValue := stock.dataClose[len(stock.dataClose)-1]

	//注意这里是状态转换：波动不管，不得不动的时候，状态才变换
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
	return ret, revert, modify, arr
}
