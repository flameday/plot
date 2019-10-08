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
		s += fmt.Sprintf(" (%d, %.2f)->(%d, %.2f) b:%.3f s:%.3f",
			int(ac.Buy_stop.left),
			ac.Buy_stop.top,
			int(ac.Buy_stop.right),
			ac.Buy_stop.bottom,
			ac.buy,
			ac.sell)
	} else if ac.Action == ACTION_SELL {
		s += fmt.Sprintf(" (%d, %.2f)->(%d, %.2f) b:%.3f s:%.3f",
			int(ac.Sell_stop.left),
			ac.Sell_stop.top,
			int(ac.Sell_stop.right),
			ac.Sell_stop.bottom,
			ac.buy,
			ac.sell)
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
		ac.profit = 0

		return true, true, true
	} else {

		ac.State = STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Sell_stop = arr[0]

		ac.sell = stock.dataClose[curIndex]
		ac.profit = 0

		return true, true, true
	}

	return true, false, false
}

func restrictStop(ac *avgContext, arr []Rect) (ret bool, revert bool, modify bool) {
	//latest 3
	size := len(arr)
	if size < 3 {
		return false, false, false
	}
	//sell
	if arr[size-1].leftFlag == 1 &&
		arr[size-1].top < arr[size-3].top+arr[size-3].Height()*0.1 &&
		arr[size-1].bottom < arr[size-3].bottom {

		change := ac.Action == STATE_NEW_LOW
		revert := ac.State == ACTION_SELL
		if ac.Action != ACTION_SELL {
			ac.sell = arr[size-1].bottom
			ac.profit += ac.sell - ac.buy
		}

		ac.State = STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Sell_stop = arr[size-1]

		return true, revert, change
	}
	//buy
	if arr[size-1].leftFlag == -1 &&
		arr[size-1].top > arr[size-3].top &&
		arr[size-1].bottom >= arr[size-3].bottom-arr[size-3].Height()*0.1 {

		change := ac.Action == STATE_NEW_HIGH
		revert := ac.State == ACTION_BUY
		if ac.Action != ACTION_BUY {
			ac.buy = arr[size-1].top
			ac.profit += ac.buy - ac.sell
		}

		ac.State = STATE_NEW_HIGH
		ac.Action = ACTION_BUY
		ac.Buy_stop = arr[size-1]

		return true, revert, change
	}

	return false, false, false
}

func action_High_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// quick change
	size := len(arr)
	ret, revert, modify = restrictStop(ac, arr)
	if ret {
		return ret, revert, modify
	}

	// keep rising top
	if arr[size-1].top > ac.Buy_stop.top {
		ac.Buy_stop.top = arr[size-1].top

		return true, false, true
	}

	// normal change
	if curValue < ac.Buy_stop.bottom {
		//new low
		if arr[size-1].top < ac.Buy_stop.top {
			ac.State = STATE_NEW_LOW
			ac.Action = ACTION_SELL
			ac.Sell_stop = ac.Buy_stop

			ac.sell = curValue
			ac.profit += ac.sell - ac.buy
		} else {
			//new high new low
			ac.State = STATE_NEW_HIGH__NEW_LOW_0
			ac.Action = ACTION_SELL
			ac.Sell_stop = ac.Buy_stop

			ac.sell = curValue
			ac.profit += ac.sell - ac.buy
		}

		return true, true, true
	}

	return true, false, false
}

func action_Low_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// quick change
	size := len(arr)
	ret, revert, modify = restrictStop(ac, arr)
	if ret {
		return ret, revert, modify
	}

	// keep fall bottom
	if arr[size-1].bottom < ac.Sell_stop.bottom {
		ac.Sell_stop.bottom = arr[size-1].bottom

		return true, false, true
	}
	// new low new high
	if curValue > ac.Sell_stop.top {
		ac.State = STATE_NEW_LOW__NEW_HIGH_0
		ac.Action = ACTION_BUY
		ac.Buy_stop = ac.Sell_stop
		ac.tmpTop = ac.Buy_stop.top

		ac.buy = curValue
		ac.profit += ac.sell - ac.buy

		return true, true, true
	}

	return true, false, false
}

func action_Low_High_0_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// quick change
	size := len(arr)
	ret, revert, modify = restrictStop(ac, arr)
	if ret {
		return ret, revert, modify
	}

	// keep rising top
	if arr[size-1].top > ac.Buy_stop.top {
		ac.Buy_stop.top = arr[size-1].top
	}

	// normal change
	if curValue < ac.Buy_stop.bottom {
		// too low
		// "A"
		ac.State = STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Buy_stop.bottom = ac.Sell_stop.bottom

		ac.sell = curValue
		ac.profit += ac.sell - ac.buy

		return true, true, true
	} else if curValue < ac.tmpTop {
		//"BC"
		ac.State = STATE_NEW_LOW__NEW_HIGH_1
		ac.Action = ACTION_SELL
		ac.Sell_stop = ac.Buy_stop

		ac.sell = curValue
		ac.profit += ac.sell - ac.buy

		return true, true, true
	}

	return true, false, false
}

func action_Low_High_1_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// quick change
	size := len(arr)
	ret, revert, modify = restrictStop(ac, arr)
	if ret {
		return ret, revert, modify
	}

	// keep fall bottom
	if arr[size-1].bottom < ac.Sell_stop.bottom {
		ac.Sell_stop.bottom = arr[size-1].bottom
	}

	if curValue < ac.tmpBottom {
		//"A"
		ac.State = STATE_NEW_LOW
		//keep sell & keep sell stop

		return true, false, true
	} else if curValue > ac.Sell_stop.top {
		//"BC"
		//in fact, no "C"
		ac.State = STATE_NEW_HIGH
		ac.Action = ACTION_BUY
		ac.Buy_stop = ac.Sell_stop

		ac.buy = curValue
		ac.profit += ac.sell - ac.buy

		return true, true, true
	}

	return true, false, false
}

func action_High_Low_0_Sell(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// quick change
	size := len(arr)
	ret, revert, modify = restrictStop(ac, arr)
	if ret {
		return ret, revert, modify
	}

	// keep fall bottom
	if arr[size-1].bottom < ac.Sell_stop.bottom {
		ac.Sell_stop.bottom = arr[size-1].bottom
	}

	// normal change
	if curValue > ac.Sell_stop.top {
		// too high
		// "A"
		ac.State = STATE_NEW_HIGH
		ac.Action = ACTION_BUY
		ac.Buy_stop = ac.Sell_stop

		ac.buy = curValue
		ac.profit += ac.sell - curValue

		return true, true, true

	} else if curValue > ac.tmpBottom {
		//"BC"
		// in fact, no "C"
		ac.State = STATE_NEW_HIGH__NEW_LOW_1
		ac.Action = ACTION_BUY
		ac.Buy_stop.top = ac.Sell_stop.top

		ac.buy = curValue
		ac.profit += ac.sell - curValue

		return true, true, true
	}

	return true, false, false
}

func action_High_Low_1_Buy(ac *avgContext, arr []Rect, curValue float64) (ret bool, revert bool, modify bool) {
	// quick change
	size := len(arr)
	ret, revert, modify = restrictStop(ac, arr)
	if ret {
		return ret, revert, modify
	}

	// keep rising top
	if arr[size-1].top > ac.Buy_stop.top {
		ac.Buy_stop.top = arr[size-1].top
	}

	// normal change
	if curValue > ac.Sell_stop.top {
		// too high
		// "A"
		ac.State = STATE_NEW_HIGH
		// keep buy & keep stop

		return true, false, false
	} else if curValue < ac.Sell_stop.bottom {
		//"BC"
		// in fact, no "C"
		ac.State = STATE_NEW_LOW
		ac.Action = ACTION_SELL
		ac.Sell_stop = ac.Buy_stop

		ac.sell = curValue
		ac.profit += ac.sell - ac.buy

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
			//p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
		}
		drawPoint(p, float64(curPos), stock.dataClose[curPos], 20, red)
		if ac.Action == ACTION_BUY {
			drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
		} else if ac.Action == ACTION_SELL {
			drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
		}
		p.X.Label.Text = ac.State + "          " + ac.Action
	} else if ac.State != STATE_UNKOWN {
		ok, revert, _, arr := forwardState(ac, stock)
		for _, r := range arr {
			drawRectangle(p, r.left, r.top, r.right, r.bottom, olive)
			if r.leftFlag == -1 {
				drawPoint(p, r.left, r.top, 10, red)
			}
			if r.leftFlag == 1 {
				drawPoint(p, r.left, r.top, 15, blue)
			}
		}

		p.X.Label.Text = ac.State + "          " + ac.Action
		drawPoint(p, float64(curPos), stock.dataClose[curPos], 20, black)
		if ac.Action == ACTION_BUY {
			drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
		} else if ac.Action == ACTION_SELL {
			drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
		}
		//log.Infof("pos:%d ok:%v", pos, ok)
		if ok && (revert) {
			if last_buy == ac.buy && last_sell == ac.sell {
			} else {
				//log.Infof("[%d] buy: %f sell:%f lb:%f, ls:%f profit:%.3f %s", pos, ac.buy, ac.sell, last_buy, last_sell, ac.profit, ac.Show())

				last_buy = ac.buy
				last_sell = ac.sell
			}
			p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
		}
	} else {
		log.Infof("ignore:%d", pos)
	}
	return true
}
