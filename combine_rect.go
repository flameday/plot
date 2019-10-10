package main

import (
	"math"
)

func GetAllRect(stock *Stock) ([]Rect, *Stock) {
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

		//if post > pre+33 {
		//	if i == 114 {
		//		log.Infof("%d --> %d diff:%d", pre, post, post-pre)
		//	}
		//	getSmallWave(stock, pre, post)
		//}

		// next
		i = post + 1
	}

	rectArray := make([]Rect, 0)
	for i := 0; i < len(stock.dataClose); {
		pre, _ := findPreMinOrMaxIndex(stock.dataMinMax, i-1)
		post := findPostMinOrMaxIndex(stock.dataMinMax, i+1)
		if post == -1 {
			i++
			continue
		}
		if pre == -1 {
			//这里加上坐标0
			if post > 0 {
				if stock.dataMinMax[post] == -1 {
					stock.dataMinMax[0] = 1
				} else {
					stock.dataMinMax[0] = -1
				}
			}

			i++
			continue
		}

		left := float64(pre)
		top := math.Max(stock.dataHigh[pre], stock.dataHigh[post])
		right := float64(post)
		bottom := math.Min(stock.dataLow[pre], stock.dataLow[post])

		offset := (top - bottom) * 0 //.05
		//offset := 0.0
		r := Rect{
			left:      left,
			top:       top + offset,
			right:     right,
			bottom:    bottom - offset,
			FlagLeft:  stock.dataMinMax[pre],
			FlagRight: stock.dataMinMax[post],
			DistLeft:  len(stock.dataClose) - pre,
			DistRight: len(stock.dataClose) - post,
		}
		rectArray = append(rectArray, r)
		i = post + 1
	}
	return rectArray, stock
}

func ExpandRect(r1 *Rect, r2 *Rect) *Rect {
	if r1.top != r2.top && r1.bottom != r1.bottom {
		panic("bad rect")
	}

	left := math.Min(r1.left, r2.left)
	right := math.Max(r1.right, r2.right)
	top := math.Max(r1.top, r2.top)
	bottom := math.Min(r1.bottom, r2.bottom)
	FlagLeft := -1
	FlagRight := 1
	if r1.top > r2.top {
		FlagLeft = 1
		FlagRight = -1
	} else if r1.bottom < r2.bottom {
		FlagLeft = 1
		FlagRight = -1
	}

	r := &Rect{
		left:      left,
		top:       top,
		right:     right,
		bottom:    bottom,
		FlagLeft:  FlagLeft,
		FlagRight: FlagRight,
	}
	return r
}
