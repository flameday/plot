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
	//flagArea    []int

	cleanPosMinMax []int
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
			//preMax := findPreIndex(stock.dataMinMax, i-1, 1)
			//postMax := findNextIndex(stock.dataMinMax, i-1, 1)
			//if preMax > preMin {
			//	continue
			//}
			//if postMax < postMin {
			//	//continue
			//}
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
			//preMin := findPreIndex(stock.dataMinMax, i-1, -1)
			//postMin := findNextIndex(stock.dataMinMax, i-1, -1)
			preMax := findPreIndex(stock.dataMinMax, i-1, 1)
			postMax := findNextIndex(stock.dataMinMax, i-1, 1)
			//if preMin > preMax {
			//	continue
			//}
			//if postMin < postMax {
			//	continue
			//}
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

func (stock *Stock) LoadData(filename string, left int, right int) (bool, int, int) {
	// 读文本数据
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Println(b)
	str := string(b)

	// 解析文本数据
	lines := strings.Split(str, "\n")
	if left >= len(lines) || right >= len(lines) {
		return false, 0, 0
	}
	if left < 0 {
		left = 0
	}
	if right > len(lines) {
		right = len(lines)
	}

	for i, line := range lines {
		if i == 0 {
			continue
		}
		// log.Info("[%d] %s", i, line)
		elems := strings.Split(line, ",")
		if len(elems) < 4 {
			continue
		}

		// 解析数值
		valueOpen, err := strconv.ParseFloat(elems[1], 32)
		if err != nil {
			log.Error("ParseFloat elems[1]:%s err:%v", elems[1], err)
		}
		valueClose, err := strconv.ParseFloat(elems[3], 32)
		if err != nil {
			log.Error("ParseFloat elems[3]:%s err:%v", elems[3], err)
		}
		//
		if i >= left && i < right {
			//if i < 500 {
			stock.dataOpen = append(stock.dataOpen, valueOpen)
			stock.dataClose = append(stock.dataClose, valueClose)
		}
	}
	//for i := 0; i < 10; i++ {
	//	log.Infof("[%d] %f", i, stock.dataClose[i])
	//}

	// 构造数据
	//stock.dataOpen = make([]float64, 0)
	//stock.dataClose = make([]float64, 0)
	//
	//for i := 1; i <= 10; i++ {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 10; i > 5; i-- {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 5; i > 15; i++ {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 15; i > 10; i-- {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 10; i > 20; i-- {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 20; i > 5; i-- {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//////
	//
	//for i := 15; i < 25; i++ {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 25; i > 15; i-- {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 15; i > 20; i++ {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 20; i < 10; i-- {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 10; i < 15; i++ {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//get_area_rate()

	//计算平均值
	for i := 0; i < len(stock.dataClose); i++ {
		val := get_pre_avg(stock.dataClose, i, 6)
		stock.avg6 = append(stock.avg6, val)

		val = get_pre_avg(stock.dataClose, i, 30)
		stock.avg30 = append(stock.avg30, val)

		val = get_pre_avg(stock.dataClose, i, 150)
		stock.avg150 = append(stock.avg150, val)

		val = get_middle_avg(stock.dataClose, i, 30)
		stock.avgMiddle = append(stock.avgMiddle, val)
	}
	//
	log.Infof("dataClose size:%d", len(stock.dataClose))
	log.Infof("avg size:%d", len(stock.avg150))
	//局部最大值
	caculateMinMax(stock.dataClose, &stock.dataMinMax, 8)
	//根据1：1的关系，过滤掉多余的大小值
	filter_max(stock.dataClose, stock.dataMinMax)
	filter_min(stock.dataClose, stock.dataMinMax)

	caculateMinMax(stock.avg30, &stock.avg30MinMax, 30)
	caculateMinMax(stock.avg150, &stock.avg150MinMax, 150)
	caculateMinMax(stock.avgMiddle, &stock.avgMiddleMinMax, 30)
	//先使用平均值的minMax，后调整
	caculateMinMax(stock.avg30, &stock.resetMinMax, 30)
	//根据平均值的大小值，往前后找真实的大小值
	locateMax(stock.dataClose, stock.resetMinMax, 61)
	locateMin(stock.dataClose, stock.resetMinMax, 61)
	//根据1：1的关系，过滤掉多余的大小值
	filter_max(stock.dataClose, stock.resetMinMax)
	filter_min(stock.dataClose, stock.resetMinMax)
	////初始化area分布
	//for i:= 0; i < len(stock.dataClose); i++{
	//	stock.flagArea = append(stock.flagArea, 0)
	//}
	//获取区间（层次）
	//locate_realate(stock.dataClose, &stock.relateCntArray)
	//run_all_caculate()

	return true, left, right
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
