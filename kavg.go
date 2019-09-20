package main

import (
	"github.com/prometheus/common/log"
)

// sequence start from 0,1,2...
func getLastHigh(data []float64, minMax []int, pos int, sequence int) (bool, float64) {
	i := pos
	for i = pos - 1; i >= 0; i-- {
		if minMax[i] == MAX_VALUE_FLAG {
			//compare
			sequence--
			if sequence <= 0 {
				return true, data[i]
			}
		}
	}
	return false, -1
}

// sequence start from 0,1,2...
func getLastLow(data []float64, minMax []int, pos int, sequence int) (bool, float64) {
	i := pos
	for i = pos - 1; i >= 0; i-- {
		if minMax[i] == MIN_VALUE_FLAG {
			//compare
			sequence--
			if sequence <= 0 {
				return true, data[i]
			}
		}
	}
	return false, -1
}

func isHigherThanLast(data []float64, minMax []int, pos int, sequence int) bool {
	curVal := data[pos]
	ok, high := getLastHigh(data, minMax, pos, sequence)
	if ok && curVal > high {
		return true
	}

	return false
}

func isLowerThanLast(data []float64, minMax []int, pos int, sequence int) bool {
	curVal := data[pos]
	ok, low := getLastLow(data, minMax, pos, sequence)
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
func actionBuy(state *int) {
	log.Infof("buy...")
	*state = STATE_BUY_FLAG
}
func actionClose(state *int) {
	log.Infof("close...")
	*state = STATE_WAIT_FLAG
}
func actionSell(state *int) {
	log.Infof("sell...")
	*state = STATE_SELL_FLAG
}

func runBuyTrend(state *int, data []float64, minMax []int, pos int) {
	ok1, firstHigh := getLastHigh(data, minMax, pos, 0)
	ok2, firstLow := getLastLow(data, minMax, pos, 0)
	ok3, secondHigh := getLastHigh(data, minMax, pos, 1)
	ok4, secondLow := getLastLow(data, minMax, pos, 1)
	if !ok1 || !ok2 {
		return
	}
	//   .    .
	//  / \  / \
	// /   .    ?
	//  移动
	if firstHigh > secondHigh {
		if !ok3 {
			return
		}
		buy_stop = firstLow
	}

	// \   .    ?
	//  \ /  \ /
	//   .    .
	if data[pos] < firstLow {
		if !ok4 {
			return
		}
		actionClose(state)
		actionSell(state)
	}
}
func runSellTrend(state *int, data []float64, minMax []int, pos int) {

}
func startTrend(state *int, data []float64, minMax []int, pos int) {
	ok1, firstHigh := getLastHigh(data, minMax, pos, 0)
	ok2, firstLow := getLastLow(data, minMax, pos, 0)
	ok3, secondHigh := getLastHigh(data, minMax, pos, 1)
	ok4, secondLow := getLastLow(data, minMax, pos, 1)
	if !ok1 || !ok2 {
		return
	}
	//   .    .
	//  / \  / \
	// /   .    ?
	if firstHigh > firstLow {
		if !ok3 {
			return
		}
		if firstHigh <= secondHigh && data[pos] < firstLow {
			actionSell(state)
		}
	}

	// \   .    ?
	//  \ /  \ /
	//   .    .
	if firstLow > firstHigh {
		if !ok4 {
			return
		}
		if firstLow <= secondLow && data[pos] > firstHigh {
			actionBuy(state)
		}
	}
}

func changeAction(state *int, data []float64, minMax []int, pos int) {
	if *state == STATE_WAIT_FLAG {
		startTrend(state, data, minMax, pos)
	} else if *state == STATE_BUY_FLAG {
		runBuyTrend(state, data, minMax, pos)
	} else if *state == STATE_SELL_FLAG {
		runSellTrend(state, data, minMax, pos)
	}
}
