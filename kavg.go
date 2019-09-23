package main

//
//import (
//	"github.com/prometheus/common/log"
//)
//
//// sequence start from 0,1,2...
//func getLastHigh(data []float64, minMax []int, pos int, sequence int) (bool, float64) {
//	i := pos
//	for i = pos - 1; i >= 0; i-- {
//		if minMax[i] == MAX_VALUE_FLAG {
//			//compare
//			sequence--
//			if sequence <= 0 {
//				return true, data[i]
//			}
//		}
//	}
//	return false, -1
//}
//
//// sequence start from 0,1,2...
//func getLastLow(data []float64, minMax []int, pos int, sequence int) (bool, float64) {
//	i := pos
//	for i = pos - 1; i >= 0; i-- {
//		if minMax[i] == MIN_VALUE_FLAG {
//			//compare
//			sequence--
//			if sequence <= 0 {
//				return true, data[i]
//			}
//		}
//	}
//	return false, -1
//}
//
//func isHigherThanLast(data []float64, minMax []int, pos int, sequence int) bool {
//	curVal := data[pos]
//	ok, high := getLastHigh(data, minMax, pos, sequence)
//	if ok && curVal > high {
//		return true
//	}
//
//	return false
//}
//
//func isLowerThanLast(data []float64, minMax []int, pos int, sequence int) bool {
//	curVal := data[pos]
//	ok, low := getLastLow(data, minMax, pos, sequence)
//	if ok && curVal < low {
//		return true
//	}
//	return false
//}
//func moveBuyStop(lastLow float64) {
//	buy_stop = lastLow
//}
//func moveSellStop(lastHigh float64) {
//	sell_stop = lastHigh
//}
//func actionBuy(state *int) {
//	log.Infof("buy...")
//	*state = ACTION_BUY_FLAG
//}
//func actionClose(state *int) {
//	log.Infof("close...")
//	//*state = ACTION_INIT_FLAG
//}
//func actionSell(state *int) {
//	log.Infof("sell...")
//	*state = ACTION_SELL_FLAG
//}
//
//func runBuyTrend(state *int, data []float64, minMax []int, pos int) {
//	ok1, firstHigh := getLastHigh(data, minMax, pos, 0)
//	ok2, firstLow := getLastLow(data, minMax, pos, 0)
//	ok3, secondHigh := getLastHigh(data, minMax, pos, 1)
//	//ok4, secondLow := getLastLow(data, minMax, pos, 1)
//	if !ok1 || !ok2 {
//		return
//	}
//	//   .    .
//	//  / \  / \
//	// /   .    ?
//	//
//	if firstLow > firstHigh {
//		//还在继续上升 或者新高又新低
//		if data[pos] > firstHigh {
//			//还在继续上升
//		}
//		return
//	}
//
//	//掉头向下了
//	if ok3 {
//		//新高
//		if firstHigh >= secondHigh {
//			//新高不新低
//			if data[pos] >= firstLow {
//				buy_stop = firstLow
//			} else if data[pos] < firstLow { //新高又新低
//				actionClose(state)
//				actionSell(state)
//				sell_stop = firstHigh
//			}
//		} else if firstHigh < secondHigh { //不新高
//			//不新高也不新低
//			if data[pos] >= firstLow {
//				//维持原状
//			} else if data[pos] < firstLow { //不新高又新低
//				actionClose(state)
//				actionSell(state)
//				sell_stop = firstHigh
//			}
//		}
//	}
//
//	// \   .    ?
//	//  \ /  \ /
//	//   .    .
//	if data[pos] < firstLow {
//		if !ok4 {
//			return
//		}
//		actionClose(state)
//		actionSell(state)
//	}
//}
//func runSellTrend(state *int, data []float64, minMax []int, pos int) {
//
//}
//func startTrend(state *int, data []float64, minMax []int, pos int) {
//	ok1, firstHigh := getLastHigh(data, minMax, pos, 0)
//	ok2, firstLow := getLastLow(data, minMax, pos, 0)
//	ok3, secondHigh := getLastHigh(data, minMax, pos, 1)
//	ok4, secondLow := getLastLow(data, minMax, pos, 1)
//	if !ok1 || !ok2 {
//		return
//	}
//	//   .    .
//	//  / \  / \
//	// /   .    ?
//	if ok3 {
//		if firstHigh > firstLow {
//			if firstHigh <= secondHigh && data[pos] < firstLow {
//				actionSell(state)
//			}
//		}
//	}
//
//	// \   .    ?
//	//  \ /  \ /
//	//   .    .
//	if ok4 {
//		if firstLow > firstHigh {
//			if firstLow <= secondLow && data[pos] > firstHigh {
//				actionBuy(state)
//			}
//		}
//	}
//}
//
//func changeAction(state *int, data []float64, minMax []int, pos int) {
//	if *state == ACTION_INIT_FLAG {
//		startTrend(state, data, minMax, pos)
//	} else if *state == ACTION_BUY_FLAG {
//		runBuyTrend(state, data, minMax, pos)
//	} else if *state == ACTION_SELL_FLAG {
//		runSellTrend(state, data, minMax, pos)
//	}
//}
