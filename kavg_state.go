package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

type Rect struct {
	left      float64
	top       float64
	right     float64
	bottom    float64
	leftFlag  int
	rightFlag int
}

func (r *Rect) isRising() bool {
	if r.left == -1 && r.right == 1 {
		return true
	} else if r.left == 1 && r.right == -1 {
		return false
	} else {
		log.Errorf("bad value")
		panic("bad value")
	}
	return false
}

func (r *Rect) Height() float64 {
	return r.top - r.bottom
}

func (r *Rect) Width() float64 {
	return r.right - r.left
}

type avgContext struct {
	State     string
	Action    string
	Sell_stop Rect
	Buy_stop  Rect
	profit    float64
	buy       float64
	sell      float64
	tmpTop    float64
	tmpBottom float64
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
	arr, _ := GetAllRect(stock)
	if len(arr) <= 0 {
		return false, false, false
	}
	curIndex := len(stock.dataClose) - 1

	//简单根据高低点来尽快买入、卖出: 注意MinMax的计算方法
	preMin := findPreIndex(stock.dataMinMax, curIndex-1, -1)
	preMax := findPreIndex(stock.dataMinMax, curIndex-1, 1)
	if preMin == -1 || preMax == -1 {
		return false, false, false
	}
	if preMin < preMax {
		ac.State = STATE_NEW_HIGH
		ac.Action = ACTION_BUY
		ac.Buy_stop = arr[0]
		ac.buy = stock.dataClose[curIndex]

		return true, true, true
	} else {

		ac.State = STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Sell_stop = arr[0]
		ac.sell = stock.dataClose[curIndex]

		return true, true, true
	}

	return true, false, false
}
func isDownWave(r Rect) bool {
	return !isRisingWave(r)
}
func isRisingWave(r Rect) bool {
	if r.leftFlag == -1 && r.rightFlag == 1 {
		return true
	} else if r.leftFlag == 1 && r.rightFlag == -1 {
		return false
	}
	log.Errorf("bad wave")
	return false
}

func action_High_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	size := len(arr)
	//如果有3个区间，就要判断 M
	if size >= 3 {
		//M
		if arr[size-3].top == ac.Sell_stop.top && arr[size-3].bottom == ac.Sell_stop.bottom {
			if curValue < arr[size-1].bottom && arr[size-1].top < arr[size-2].top+arr[size-2].Width()*0.1 {
				ac.State = STATE_NEW_LOW
				ac.Action = ACTION_SELL
				ac.Sell_stop = arr[size-1]
				ac.sell = curValue
				return true, true, true
			}
		}
	}

	// 正常 W
	if curValue >= ac.Buy_stop.top && size >= 3 {
		// 如果是上涨区间，推进 Buy_stop
		if arr[size-3].top == ac.Sell_stop.top && arr[size-3].bottom == ac.Sell_stop.bottom {
			if curValue > arr[size-2].top && arr[size-2].bottom > arr[size-3].bottom-arr[size-3].Height()*0.1 {
				ac.State = STATE_NEW_HIGH
				ac.Action = ACTION_BUY
				ac.Buy_stop = arr[size-1]
				ac.buy = curValue
				return true, true, true
			}
		}
	} else if curValue < ac.Buy_stop.bottom {
		//必定是下跌
		if arr[size-1].leftFlag == -1 {
			log.Errorf("bad rect")
		}
		//新低
		if arr[size-1].top < ac.Buy_stop.top {
			ac.State = STATE_NEW_LOW
			ac.Action = ACTION_SELL
			ac.Sell_stop = ac.Buy_stop
			ac.sell = curValue
			ac.profit += curValue - ac.buy
		} else {
			//新高又新低
			ac.State = STATE_NEW_HIGH__NEW_LOW_0
			ac.Action = ACTION_SELL
			ac.Sell_stop = ac.Buy_stop

			ac.sell = curValue
			ac.profit += curValue - ac.buy
		}

		return true, true, true
	}

	if arr[size-1].top > ac.Buy_stop.top {
		ac.Buy_stop.top = arr[size-1].top

		return true, false, true
	}

	return true, false, false
}
func restrictStop(ac *avgContext, arr []Rect) (ret bool, revert bool, modify bool) {
	//按最近的3个判断
	size := len(arr)
	if size < 3 {
		return false, false, false
	}
	//sell
	if arr[size-1].top < arr[size-3].top+arr[size-3].Height()*0.1 && arr[size-1].bottom < arr[size-3].bottom {
		ac.State = STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Buy_stop = arr[size-1]
		return true, true, true
	}
	//buy
	if arr[size-1].top > arr[size-3].top && arr[size-1].bottom > arr[size-3].bottom {
		ac.State = STATE_NEW_HIGH
		ac.Action = ACTION_BUY
		ac.Buy_stop = arr[size-1]
		return true, true, true
	}

	return false, false, false
}

func action_Low_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	size := len(arr)
	//如果有3个区间，就要判断 W
	if size >= 4 {
		//W
		if arr[size-4].top == ac.Sell_stop.top && arr[size-4].bottom == ac.Sell_stop.bottom {
			if curValue > arr[size-2].top && arr[size-2].bottom > arr[size-3].bottom-arr[size-3].Height()*0.1 {
				ac.State = STATE_NEW_HIGH
				ac.Action = ACTION_BUY
				ac.Buy_stop = arr[size-1]
				ac.buy = curValue
				return true, true, true
			}
		}
	}

	// 正常 M
	if curValue <= ac.Sell_stop.top && size >= 3 {
		// 如果是下跌区间，推进 Sell_stop
		//M
		if arr[size-3].top == ac.Sell_stop.top && arr[size-3].bottom == ac.Sell_stop.bottom {
			if curValue < arr[size-2].bottom && arr[size-1].top < arr[size-2].top+arr[size-2].Height()*0.1 {
				ac.State = STATE_NEW_LOW
				ac.Action = ACTION_SELL
				ac.Sell_stop = arr[size-1]
				ac.sell = curValue
				return true, true, true
			}
		}
	}
	// 新低
	if curValue > ac.Sell_stop.top {
		ac.State = STATE_NEW_LOW__NEW_HIGH_0
		ac.Action = ACTION_BUY
		ac.Buy_stop = ac.Sell_stop
		ac.tmpTop = ac.Buy_stop.top

		ac.buy = curValue
		ac.profit += ac.sell - curValue

		return true, true, true
	}

	//bottom 降低
	if arr[size-1].bottom < ac.Sell_stop.bottom {
		ac.Sell_stop.bottom = arr[size-1].bottom

		return true, false, true
	}

	return true, false, false
}
func action_Low_High_0_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {

	size := len(arr)
	if arr[size-1].top > ac.Buy_stop.top {
		ac.Buy_stop.top = arr[size-1].top
	}

	if curValue < ac.Buy_stop.bottom {
		//"A"
		ac.State = STATE_NEW_LOW
		// 特殊情况，一步跳到位
		if ac.Action == ACTION_BUY {
			ac.Action = ACTION_SELL
			ac.Buy_stop.bottom = ac.Sell_stop.bottom

			ac.sell = curValue
			ac.profit += curValue - ac.buy

			return true, true, true
		}
		return true, false, false
	} else if curValue < ac.tmpTop {
		//"BC"
		ac.State = STATE_NEW_LOW__NEW_HIGH_1
		ac.Action = ACTION_SELL
		ac.Sell_stop = ac.Buy_stop

		ac.sell = curValue
		ac.profit += curValue - ac.buy

		return true, true, true
	}

	return true, false, false
}

func action_Low_High_1_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	size := len(arr)
	//如果有3个区间，就要判断 M
	if size >= 4 {
		//M
		if arr[size-4].top == ac.Sell_stop.top && arr[size-4].bottom == ac.Sell_stop.bottom {
			if curValue < arr[size-2].bottom && arr[size-1].top < arr[size-3].top+arr[size-3].Width()*0.1 {
				ac.State = STATE_NEW_LOW
				ac.Action = ACTION_SELL
				ac.Sell_stop = arr[size-1]
				ac.sell = curValue
				return true, true, true
			}
		}
	}

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

		return true, true, true
	}

	return true, false, false
}

func action_High_Low_0_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// 更新最小值
	if curValue > ac.Sell_stop.top {
		// "A"
		// 这个得排在前面，优先级更高
		ac.State = STATE_NEW_HIGH
		// 特殊情况，一步跳到位
		if ac.Action == ACTION_SELL {
			ac.Action = ACTION_BUY
			ac.Buy_stop.top = ac.Sell_stop.top

			ac.buy = curValue
			ac.profit += ac.sell - curValue

			return true, true, true
		}

		return true, false, true
	} else if curValue > ac.Sell_stop.bottom {
		//"BC"
		// in fact, no "C"
		ac.State = STATE_NEW_HIGH__NEW_LOW_1
		ac.Action = ACTION_BUY
		ac.Buy_stop.top = ac.Sell_stop.top

		ac.buy = curValue
		ac.profit += ac.sell - curValue

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

		return true, true, true
	}

	return true, false, false
}

func forwardState(ac *avgContext, stock *Stock) (ret bool, revert bool, modify bool, arr []Rect) {
	//get pre arr
	arr, _ = GetAllRect(stock)
	if len(arr) <= 0 {
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

func Run(ac *avgContext, p *plot.Plot, stock *Stock, filename string, pos int) bool {

	if pos == 350 {
		log.Errorf("350")
	}
	//_, st := GetAllRect(stock)
	curPos := len(stock.dataClose) - 1

	if ac.State == STATE_UNKOWN {
		ok, revert, change := isValidInit(ac, stock)
		if ok && revert && change {
			//log.Infof("ac: %s", ac.Show())
			p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
		}
		drawPoint(p, float64(curPos), stock.dataClose[curPos], 20, red)
		if ac.Action == ACTION_BUY {
			drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
		} else if ac.Action == ACTION_SELL {
			drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
		}
		p.X.Label.Text = ac.State + "          " + ac.Action
	} else if ac.State != STATE_UNKOWN {
		ok, revert, change, arr := forwardState(ac, stock)
		for _, r := range arr {
			drawRectangle(p, r.left, r.top, r.right, r.bottom, olive)
			if r.leftFlag == -1 {
				drawPoint(p, r.left, r.top, 10, red)
			}
			if r.leftFlag == 1 {
				drawPoint(p, r.left, r.top, 15, blue)
			}
		}
		//log.Infof("pos:%d ok:%v", pos, ok)
		if ok && (revert || change) {
		}

		p.X.Label.Text = ac.State + "          " + ac.Action
		drawPoint(p, float64(curPos), stock.dataClose[curPos], 20, black)
		if ac.Action == ACTION_BUY {
			drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
		} else if ac.Action == ACTION_SELL {
			drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
		}
		//p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
	} else {
		log.Infof("ignore:%d", pos)
	}
	return false
}
