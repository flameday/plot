package main

import (
	"github.com/prometheus/common/log"
)

func getLastHigh(data []float64, minMax []int, pos int) (bool, float64) {
	i := pos
	for i = pos - 1; i >= 0; i-- {
		if minMax[i] == MAX_VALUE_FLAG {
			//compare
			return true, data[i]
		}
	}
	return false, -1
}

func getLastLow(data []float64, minMax []int, pos int) (bool, float64) {
	i := pos
	for i = pos - 1; i >= 0; i-- {
		if minMax[i] == MIN_VALUE_FLAG {
			//compare
			return true, data[i]
		}
	}
	return false, -1
}

func isHigherThanLast(data []float64, minMax []int, pos int) bool {
	curVal := data[pos]
	ok, high := getLastHigh(data, minMax, pos)
	if ok && curVal > high {
		return true
	}

	return false
}

func isLowerThanLast(data []float64, minMax []int, pos int) bool {
	curVal := data[pos]
	ok, low := getLastLow(data, minMax, pos)
	if ok && curVal < low {
		return true
	}
	return false
}
func moveBuyStop(lastLow float64) {
	buy_stop = lastLow
}
func moveSellStop(lastHigh float64) {
	sell_stop = lastHigh
}
func actionBuy() {
	log.Infof("buy...")
}
func actionClose() {
	log.Infof("close...")
}
func actionSell() {
	log.Infof("sell...")
}
func changeAction(state int, data []float64, minMax []int, pos int) {
	flagHigherLast := isHigherThanLast(data, minMax, pos)
	flagLowerLast := isLowerThanLast(data, minMax, pos)

	//lastHigh := getLastHigh(data, minMax, pos)
	//lastLow := getLastLow(data, minMax, pos)
	if state == STATE_WAIT_FLAG {
		if flagHigherLast {
			actionBuy()
		} else if flagLowerLast {
			actionSell()
		}
	} else if state == STATE_BUY_FLAG {
		if flagHigherLast {
			//keep
		} else if flagLowerLast {
			actionClose()
			actionSell()
		}
	} else if state == STATE_SELL_FLAG {
		if flagHigherLast {
			actionClose()
			actionBuy()
		} else if flagLowerLast {
			//actionSell()
		}
	}
}
