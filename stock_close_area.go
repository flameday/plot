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
var puzzleAreaArray [1000][1000]float64
var puzzleBoundArray [1000][1000]float64
var puzzelRateArray [1000][1000]float64

func run_all_caculate() {
	init_low_high_diagonal()
	get_full_low_high_diagonal()
	show_low_high_diagonal()

	init_puzzle_area_diagonal()
	get_full_puzzle_area_diagonal()
	show_puzzle_area_diagonal()

	get_full_puzzle_bound_diagonal()
	show_puzzle_bound_diagonal()
}

func init_low_high_diagonal() {
	for i := 0; i < len(stock.dataClose); i++ {
		minValue := stock.dataClose[i]
		maxValue := stock.dataClose[i]
		val := lowHigh{
			Low:  minValue,
			High: maxValue,
		}
		lowHighArray[i][i] = val
	}
}

func get_full_low_high_diagonal() {
	for m := 1; m < len(stock.dataClose); m++ {
		for k := 0; k < len(stock.dataClose)-m; k++ {
			x := m + k
			y := k
			//log.Infof("low+high=(%d, %d)", x, y)

			left := lowHighArray[x-1][y]
			right := lowHighArray[x][y+1]
			val := lowHigh{
				Low:  math.Min(left.Low, right.Low),
				High: math.Max(left.High, right.High),
			}
			lowHighArray[x][y] = val
			//log.Infof("===[%d][%d] val:%.2f left:%f right:%.2f", y+x, y, val, left, right)
		}
	}
}

func show_low_high_diagonal() {
	log.Infof("-----------------------")
	tmp := "lowHigh:     "
	for i := 0; i < len(stock.dataClose); i++ {
		tmp += fmt.Sprintf("[%3d] ", i)
	}
	log.Infof(tmp)
	for y := 0; y < len(stock.dataClose); y++ {
		res := ""
		for x := 0; x < len(stock.dataClose); x++ {
			res += fmt.Sprintf("%2d/%-2d ", int(lowHighArray[x][y].Low), int(lowHighArray[x][y].High))
		}
		log.Infof("lowHigh:[%2d] %s", y, res)
	}
}

//------------------------------------------
func init_puzzle_area_diagonal() {
	for i := 0; i < len(stock.dataClose); i++ {
		puzzleAreaArray[i][i] = 0
	}
	for i := 1; i < len(stock.dataClose); i++ {
		puzzleAreaArray[i][i-1] = math.Abs(stock.dataClose[i] - stock.dataClose[i-1])
	}
}
func get_full_puzzle_area_diagonal() {
	for m := 2; m < len(stock.dataClose); m++ {
		for k := 0; k < len(stock.dataClose)-m; k++ {
			x := m + k
			y := k
			//log.Infof("x:%d y:%d", x, y)
			left := puzzleAreaArray[x-1][y]
			right := puzzleAreaArray[x][y+1]
			val := left + right
			puzzleAreaArray[x][y] = val
			//log.Infof("===[%d][%d] val:%.2f left:%f right:%.2f", y+x, y, val, left, right)
		}
	}
}
func show_puzzle_area_diagonal() {
	log.Infof("-----------------------")
	tmp := "puzzle area:     "
	for i := 0; i < len(stock.dataClose); i++ {
		tmp += fmt.Sprintf("[%2d] ", i)
	}
	log.Infof(tmp)

	for y := 0; y < len(stock.dataClose); y++ {
		res := ""
		for x := 0; x < len(stock.dataClose); x++ {
			res += fmt.Sprintf("%4d ", int(puzzleAreaArray[x][y]))
		}
		log.Infof("puzzle area:[%2d] %s", y, res)
	}
}

//------------------------------------------
func get_full_puzzle_bound_diagonal() {
	for m := 0; m < len(stock.dataClose); m++ {
		for k := 0; k < len(stock.dataClose)-m; k++ {
			x := m + k
			y := k
			diff := math.Abs(float64(x-y)) + 1
			low := lowHighArray[x][y].Low
			high := lowHighArray[x][y].High
			bound := math.Abs(high-low) + 1
			//log.Infof("bound===>(%2d %2d) %f", x, y, diff*bound)

			//if x == 8 && y == 2 {
			//	x = 8 * y / 18 / 8
			//	log.Infof("pause")
			//}
			puzzleBoundArray[x][y] = diff * bound
		}
	}
}

func show_puzzle_bound_diagonal() {
	log.Infof("-----------------------")
	tmp := "puzzle bound:     "
	for i := 0; i < len(stock.dataClose); i++ {
		tmp += fmt.Sprintf("[%4d]", i)
	}
	log.Infof(tmp)

	for y := 0; y < len(stock.dataClose); y++ {
		res := ""
		for x := 0; x < len(stock.dataClose); x++ {
			res += fmt.Sprintf("%6d", int(puzzleBoundArray[x][y]))
		}
		log.Infof("puzzle bound:[%2d] %s", y, res)
	}
}

//--------------------------------------------------------
func reset_puzzle_diagonal() {
	for i := 0; i < len(stock.dataClose); i++ {
		puzzelRateArray[i][i] = 0
	}
}
func get_full_puzzle_rate_diagnoal() {
	for m := 0; m < len(stock.dataClose); m++ {
		for k := 0; k < len(stock.dataClose)-m; k++ {
			x := m + k
			y := k
			area := puzzleAreaArray[x][y]
			bound := puzzleBoundArray[x][y]
			if x == y {
				//puzzelRateArray[x][y] = 0
				puzzelRateArray[x][y] = area / bound
			} else {
				puzzelRateArray[x][y] = area / bound
			}
		}
	}
}
