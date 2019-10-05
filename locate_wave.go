package main

import (
	log "github.com/cihub/seelog"
	"math"
)

func firstWave(stock *Stock, index int) {
	//minMax := stock.dataMinMax
	prePos, _ := findPreMinOrMaxIndex(stock.dataMinMax, index-1)
	if prePos == -1 {
		prePos = 0
	}
	highCnt := 0
	lowCnt := 0
	for i := prePos; i <= index; i++ {
		if stock.dataHigh[i] < stock.avg10[i] {
			lowCnt += 1
		} else if stock.dataLow[i] > stock.avg10[i] {
			highCnt += 1
		}
	}
	if highCnt >= 3 {
		stock.dataMinMax[index] = 1
	} else if lowCnt >= 3 {
		stock.dataMinMax[index] = -1
	}
}

func getWave(stock *Stock, index int) {
	minMax := stock.dataMinMax

	pos, _ := findPreMinOrMaxIndex(minMax, index-1)
	if pos == -1 {
		firstWave(stock, index)
		return
	}
	preLowPos := findPreIndex(minMax, index-1, -1)
	preHighPos := findPreIndex(minMax, index-1, 1)

	highCntFromLow := 0
	lowCntFromLow := 0

	highCntFromHigh := 0
	lowCntFromHigh := 0
	tmpPreLowPos := preLowPos
	if tmpPreLowPos == -1 {
		tmpPreLowPos = 0
	}
	for i := tmpPreLowPos; i <= index; i++ {
		if stock.dataHigh[i] < stock.avg10[i] {
			lowCntFromLow += 1
		} else if stock.dataLow[i] > stock.avg10[i] {
			highCntFromLow += 1
		}
	}

	tmpPreHighPos := preHighPos
	if tmpPreHighPos == -1 {
		tmpPreHighPos = 0
	}
	for i := tmpPreHighPos; i <= index; i++ {
		if stock.dataHigh[i] < stock.avg10[i] {
			lowCntFromHigh += 1
		} else if stock.dataLow[i] > stock.avg10[i] {
			highCntFromHigh += 1
		}
	}
	//可能有一个为 -1
	//低 ---> 高
	if preLowPos < preHighPos || preLowPos == -1 {
		if highCntFromLow >= 3 {
			//merge
			if stock.dataHigh[preHighPos] > stock.dataHigh[index] {
				//keep old
			} else if stock.dataHigh[index] > stock.avg10[index] {
				//use new
				minMax[preHighPos] = 0
				minMax[index] = 1
			}
		}
		if lowCntFromHigh >= 3 {
			minMax[index] = -1
		}
	}
	//高 ---> 低
	if preHighPos < preLowPos || preHighPos == -1 {
		if lowCntFromHigh >= 3 {
			//merge
			if stock.dataLow[preLowPos] < stock.dataLow[index] {
				//keep old
			} else if stock.dataLow[index] < stock.avg10[index] {
				//use new
				minMax[preLowPos] = 0
				minMax[index] = -1
			}
		}

		if highCntFromLow >= 3 {
			minMax[index] = 1
		}
	}
}

// 根据前面的 length 个值，获取平均值
func get_pre_avg(value_list []float64, index int, length int) float64 {
	var total float64
	var count int
	total = 0.0
	count = 0
	for i := 0; i < length; i++ {
		pos := index - i
		if pos < 0 {
			break
		}
		total += value_list[pos]
		count += 1
	}
	return total / float64(count)
}

func getAllRect(stock *Stock) ([]Rect, *Stock) {
	//计算最大最小值
	for i := 0; i < len(stock.dataClose); i++ {
		//log.Infof("i: %d", i)
		getWave(stock, i)
	}
	//找到最后一个点
	lastPos, _ := findPreMinOrMaxIndex(stock.dataMinMax, len(stock.dataMinMax)-1)
	lastPos, _ = findPreMinOrMaxIndex(stock.dataMinMax, lastPos-1)
	for i := 0; i < lastPos; {
		pre, _ := findPreMinOrMaxIndex(stock.dataMinMax, i-1)
		if pre == -1 {
			i++
			continue
		}
		post := findPostMinOrMaxIndex(stock.dataMinMax, i+1)
		if post == -1 {
			i++
			continue
		}

		if post > pre+33 {
			if i == 114 {
				log.Infof("%d --> %d diff:%d", pre, post, post-pre)
			}
			getSmallWave(stock, pre, post)
		}

		// next
		i = post + 1
	}

	rectArray := make([]Rect, 0)
	for i := 0; i < len(stock.dataClose); {
		pre, _ := findPreMinOrMaxIndex(stock.dataMinMax, i-1)
		if pre == -1 {
			i++
			continue
		}
		post := findPostMinOrMaxIndex(stock.dataMinMax, i+1)
		if post == -1 {
			i++
			continue
		}

		left := float64(pre)
		top := math.Max(stock.dataHigh[pre], stock.dataHigh[post])
		right := float64(post)
		bottom := math.Min(stock.dataLow[pre], stock.dataLow[post])

		offset := (top - bottom) * 0.05
		//offset := 0.0
		r := Rect{
			left:      left,
			top:       top + offset,
			right:     right,
			bottom:    bottom - offset,
			leftFlag:  stock.dataMinMax[pre],
			rightFlag: stock.dataMinMax[post],
		}
		rectArray = append(rectArray, r)
		i = post + 1
	}
	return rectArray, stock
}
