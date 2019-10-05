package main

import (
	log "github.com/cihub/seelog"
	"gonum.org/v1/plot"
)

func drawWave(p *plot.Plot, stock *Stock) {
	arr, _ := GetAllRect(stock)
	if len(arr) > 0 {
		//for _, r := range arr {
		//	drawRectangle(p, r.left, r.top, r.right, r.bottom, gray)
		//}

		//find centrum
		flagArr := mergeRect(arr)
		log.Infof("flagArr:%v", flagArr)

		lastRect := &Rect{0, 0, 0, 0, -1, -1}
		for k := 0; k < len(flagArr); k++ {
			if flagArr[k] == 0 {
				//不粘
				if lastRect.left == 0 && lastRect.right == 0 {
					lastRect = &arr[k]
				} else {
					lastRect = ExpandRect(lastRect, &arr[k])
					//draw expand
					drawRectangle2(p, lastRect.left, lastRect.top, lastRect.right, lastRect.bottom, 4, red)
				}
				lastRect = &Rect{0, 0, 0, 0, -1, -1}
			} else {
				if lastRect.left == 0 && lastRect.right == 0 {
					lastRect = &arr[k]
				} else {
					lastRect = ExpandRect(lastRect, &arr[k])
				}
			}
		}
	}
}
