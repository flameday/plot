package main

import (
	"github.com/prometheus/common/log"
)

type lowHigh struct {
	low float64
	high float64
}

func get_area_rate(dataClose []float64, index int, length int) float64 {
	if index < 0 || index + length > len(dataClose) {
		return -1
	}
	//找最高最低点

}
// 获取符合要求的个数
func get_relate_cnt(data []float64, index int, length int, xrate float64, yrate float64) float64 {
	cnt := 0.0
	for i:= index - length/2; i < index + length/2; i++ {
		if i < 0 {
			continue
		}
		if i >= len(data) {
			continue
		}
		deltaY := (data[i] - data[index])
		deltaX := i - index
		absValue := deltaY*deltaY*xrate + float64(deltaX*deltaX)*yrate
		log.Infof("        [%d] absValue:%f", i, absValue)
		if (absValue < 30) {
			cnt += 1
		}
	}
	return cnt;
}

func locate_realate(data []float64, dstArray *[]float64) {
	for i:= 0;i < len(data); i++ {
		cnt := get_relate_cnt(data, i, 20, 100, 1)
		log.Infof("[%d] cnt:%f", i, cnt)

		*dstArray = append(*dstArray, cnt/10)
	}
}
//func caculateSquare() float64 {
//
//}

//func get_area() {
//	//从第一个高点开始，找后面的低点
//	for i := 0; i < len(stock.resetMinMax); i++ {
//		pos := findNextIndex(stock.resetMinMax, i, 1)
//		if pos == -1 {
//			//找不到高点了
//			break
//		}
//
//		//往右边找点
//		rightStartPos := pos
//		validPosArr := make([]int, 0)
//		validPosArr = append(validPosArr, pos)
//		for j := pos; j < len(stock.resetMinMax); j++ {
//			rightNewPos := findNextIndex(stock.resetMinMax, rightStartPos, -1)
//			if (rightNewPos == -1) {
//				break
//			}
//			//计算面积
//
//
//			rightStartPos = rightNewPos
//
//			stock.flagArea[rightStartPos] = -1
//		}
//
//	}
//}
