package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type Stock struct {
	dataClose      []float64
	dataOpen       []float64
	avg6           []float64
	avg30          []float64
	avg150         []float64
	avgMiddle      []float64
	relateCntArray []float64

	dataMinMax      []int
	avg30MinMax     []int
	avg150MinMax    []int
	avgMiddleMinMax []int

	resetMinMax []int

	cleanPosMinMax []int
}

type Rect struct {
	left   float64
	top    float64
	right  float64
	bottom float64
}

type avgContext struct {
	State        int
	SubState     string
	Action       string
	Sell_stop    Rect
	Buy_stop     Rect
	Min_High_low float64
}

func (stock *Stock) LoadAllData(filename string) {
	// 读文本数据
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	// 解析文本数据
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		// log.Info("[%d] %s", i, line)
		elems := strings.Split(line, ",")
		if len(elems) < 4 {
			continue
		}

		valueClose, err := strconv.ParseFloat(elems[3], 32)
		if err != nil {
			log.Error("ParseFloat elems[3]:%s err:%v", elems[3], err)
		}
		stock.dataClose = append(stock.dataClose, valueClose)
	}
	log.Infof("filename: %s", filename)
	log.Infof("stock.dataClose size: %d", len(stock.dataClose))
}

func getAllRect(data []float64) ([]Rect, *Stock) {
	var stock = Stock{
		dataClose: data,
	}
	// 局部最大值
	caculateMinMax(stock.dataClose, &stock.dataMinMax, 8)
	//根据1：1的关系，过滤掉多余的大小值
	filter_max(stock.dataClose, stock.dataMinMax)
	filter_min(stock.dataClose, stock.dataMinMax)

	rectArray := make([]Rect, 0)
	// 查找
	for i := 1; i < len(stock.dataMinMax)-1; i++ {
		if stock.dataMinMax[i] == 1 {
			preMin := findPreIndex(stock.dataMinMax, i-1, -1)
			postMin := findNextIndex(stock.dataMinMax, i-1, -1)

			if preMin != -1 && postMin != -1 {
				x1, x2 := getRectangle(stock.dataClose, preMin, i, postMin)
				left := math.Min(float64(x1), float64(x2))
				top := math.Max(stock.dataClose[x1], stock.dataClose[x2])
				right := math.Max(float64(x1), float64(x2))
				bottom := math.Min(stock.dataClose[x1], stock.dataClose[x2])
				r := Rect{
					left:   left,
					top:    top,
					right:  right,
					bottom: bottom,
				}
				rectArray = append(rectArray, r)
			}
		}
		if stock.dataMinMax[i] == -1 {
			preMax := findPreIndex(stock.dataMinMax, i-1, 1)
			postMax := findNextIndex(stock.dataMinMax, i-1, 1)

			if preMax != -1 && postMax != -1 {
				x1, x2 := getRectangle(stock.dataClose, preMax, i, postMax)
				left := math.Min(float64(x1), float64(x2))
				top := math.Max(stock.dataClose[x1], stock.dataClose[x2])
				right := math.Max(float64(x1), float64(x2))
				bottom := math.Min(stock.dataClose[x1], stock.dataClose[x2])
				r := Rect{
					left:   left,
					top:    top,
					right:  right,
					bottom: bottom,
				}
				rectArray = append(rectArray, r)
			}
		}
	}
	return rectArray, &stock
}

func caculateMinMax(dataClose []float64, minMax *[]int, length int) {
	for i := 0; i < len(dataClose); i++ {
		flag := 0
		if isMax(dataClose, i, length) {
			flag = 1
		} else if isMin(dataClose, i, length) {
			flag = -1
		}
		*minMax = append(*minMax, flag)
	}
}
