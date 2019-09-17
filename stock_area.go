package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"math"
)

type lowHigh struct {
	Low  float64
	High float64
}

var lowHighArray [10000][10000]lowHigh
var candleSpaceArray [10000][10000]float64

//func get_area_rate(dataClose []float64, index int, length int) float64 {
//	if index < 0 || index+length > len(dataClose) {
//		return -1
//	}

func get_full_low_high_diagonal() {
	for i := 1; i < len(stock.dataOpen); i++ {
		for j := 0; j < len(stock.dataOpen)-i; j++ {

			left := lowHighArray[j+i-1][j]
			right := lowHighArray[j+i][j+1]
			val := lowHigh{
				Low:  math.Min(left.Low, right.Low),
				High: math.Max(left.High, right.High),
			}
			lowHighArray[j+i][j] = val
			log.Infof("===[%d][%d] val:%.2f left:%f right:%.2f", j+i, j, val, left, right)
		}
	}
}
func get_full_candle_space() {
	for i := 1; i < len(stock.dataOpen); i++ {
		for j := 0; j < len(stock.dataOpen)-i; j++ {

			left := candleSpaceArray[j+i-1][j]
			right := candleSpaceArray[j+i][j+1]
			val := left + right
			candleSpaceArray[j+i][j] = val
			log.Infof("===[%d][%d] val:%.2f left:%f right:%.2f", j+i, j, val, left, right)
		}
	}
}
func show_candle_space() {
	for j := 0; j < len(stock.dataOpen); j++ {
		res := ""
		for i := 0; i < len(stock.dataOpen); i++ {
			res += fmt.Sprintf(" %.2f", candleSpaceArray[i][j])
		}
		log.Infof("candle:[%d] %s", j, res)
	}
}
func init_low_high_diagonal() {
	for i := 0; i < len(stock.dataOpen); i++ {
		minValue := math.Min(stock.dataOpen[i], stock.dataClose[i])
		maxValue := math.Max(stock.dataOpen[i], stock.dataClose[i])
		val := lowHigh{
			Low:  minValue,
			High: maxValue,
		}
		lowHighArray[i][i] = val
	}
}
func init_diagonal() {
	for i := 0; i < len(stock.dataOpen); i++ {
		diff := math.Abs(stock.dataOpen[i]-stock.dataClose[i]) + 0.01
		candleSpaceArray[i][i] = diff
	}
}

func get_area_rate() float64 {
	//找最高最低点
	for i := 0; i < len(stock.dataOpen); i++ {
		for j := 0; j < len(stock.dataOpen)-i; j++ {
			//log.Infof("===[%d][%d]", j, j+i)
			//first := j
			//second := j + i
			//maxValue := math.Max(stock.dataOpen[first], stock.dataOpen[second])
			//maxValue := math.Max(maxValue, stock.dataClose[second])
			//maxValue := math.Max(maxValue, stock.dataClose[second])
		}
	}
	init_diagonal()
	get_full_candle_space()
	show_candle_space()

	return -1
}

// 获取符合要求的个数
func get_relate_cnt(data []float64, index int, length int, xrate float64, yrate float64) float64 {
	cnt := 0.0
	for i := index - length/2; i < index+length/2; i++ {
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
		if absValue < 30 {
			cnt += 1
		}
	}
	return cnt
}

func locate_realate(data []float64, dstArray *[]float64) {
	for i := 0; i < len(data); i++ {
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
