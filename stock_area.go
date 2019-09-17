package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"math"
)

type lowHigh struct {
	Low  float64
	High float64
}

var lowHighArray [1000][1000]lowHigh
var candleSpaceArray [1000][1000]float64
var puzzleAreaArray [1000][1000]float64
var puzzleBoundArray [1000][1000]float64

//func get_area_rate(dataClose []float64, index int, length int) float64 {
//	if index < 0 || index+length > len(dataClose) {
//		return -1
//	}

func get_full_puzzle_bound_diagonal() {
	for x := 0; x < len(stock.dataOpen); x++ {
		for y := 0; y < len(stock.dataOpen)-x; y++ {
			diff := math.Abs(float64(x-y)) + 1
			low := lowHighArray[x][y].Low
			high := lowHighArray[x][y].High
			bound := math.Abs(high-low) + 0.01

			puzzleBoundArray[x][y] = diff * bound
		}
	}
}
func get_full_puzzle_area_diagonal() {
	for x := 1; x < len(stock.dataOpen); x++ {
		for y := 0; y < len(stock.dataOpen)-x; y++ {
			left := puzzleAreaArray[y+x-1][y]
			right := puzzleAreaArray[y][y]
			val := left + right
			puzzleAreaArray[y+x][y] = val
			//log.Infof("===[%d][%d] val:%.2f left:%f right:%.2f", y+x, y, val, left, right)
		}
	}
}
func get_full_low_high_diagonal() {
	for x := 1; x < len(stock.dataOpen); x++ {
		for y := 0; y < len(stock.dataOpen)-x; y++ {

			left := lowHighArray[y+x-1][y]
			right := lowHighArray[y+x][y+1]
			val := lowHigh{
				Low:  math.Min(left.Low, right.Low),
				High: math.Max(left.High, right.High),
			}
			lowHighArray[y+x][y] = val
			//log.Infof("===[%d][%d] val:%.2f left:%f right:%.2f", y+x, y, val, left, right)
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
			//log.Infof("===[%d][%d] val:%.2f left:%f right:%.2f", j+i, j, val, left, right)
		}
	}
}

func show_puzzle_bound_diagonal() {
	log.Infof("-----------------------")
	for y := 0; y < len(stock.dataOpen); y++ {
		res := ""
		for x := 0; x < len(stock.dataOpen); x++ {
			res += fmt.Sprintf("%6d", int(puzzleBoundArray[x][y]))
		}
		log.Infof("puzzle bound:[%2d] %s", y, res)
	}
}
func show_puzzle_area_diagonal() {
	log.Infof("-----------------------")
	for y := 0; y < len(stock.dataOpen); y++ {
		res := ""
		for x := 0; x < len(stock.dataOpen); x++ {
			res += fmt.Sprintf("%4d ", int(puzzleAreaArray[x][y]))
		}
		log.Infof("puzzle area:[%2d] %s", y, res)
	}
}
func show_low_high_diagonal() {
	log.Infof("-----------------------")
	for y := 0; y < len(stock.dataOpen); y++ {
		res := ""
		for x := 0; x < len(stock.dataOpen); x++ {
			res += fmt.Sprintf("%2d/%-2d ", int(lowHighArray[x][y].Low), int(lowHighArray[x][y].High))
		}
		log.Infof("lowHigh:[%2d] %s", y, res)
	}
}
func show_candle_space() {
	log.Infof("-----------------------")
	for j := 0; j < len(stock.dataOpen); j++ {
		res := ""
		for i := 0; i < len(stock.dataOpen); i++ {
			res += fmt.Sprintf(" %.2f", candleSpaceArray[i][j])
		}
		log.Infof("candle:[%2d] %s", j, res)
	}
}

func init_puzzle_area_diagonal() {
	for i := 0; i < len(stock.dataOpen); i++ {
		puzzleAreaArray[i][i] = math.Abs(lowHighArray[i][i].High - lowHighArray[i][i].Low)
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
	//init_diagonal()
	//get_full_candle_space()
	//show_candle_space()

	init_low_high_diagonal()
	get_full_low_high_diagonal()
	show_low_high_diagonal()

	init_puzzle_area_diagonal()
	get_full_puzzle_area_diagonal()
	show_puzzle_area_diagonal()

	get_full_puzzle_bound_diagonal()
	show_puzzle_bound_diagonal()

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
		//log.Infof("        [%d] absValue:%f", i, absValue)
		if absValue < 30 {
			cnt += 1
		}
	}
	return cnt
}

func locate_realate(data []float64, dstArray *[]float64) {
	for i := 0; i < len(data); i++ {
		cnt := get_relate_cnt(data, i, 20, 100, 1)
		//log.Infof("[%d] cnt:%f", i, cnt)

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
