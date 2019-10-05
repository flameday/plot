package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"strconv"
	"strings"
)

type Stock struct {
	dataClose        []float64
	dataOpen         []float64
	dataHigh         []float64
	dataLow          []float64
	subBar           []float64
	avg10            []float64
	avg30            []float64
	avg150           []float64
	avgMiddle        []float64
	relateCntArray   []float64
	normalizedAvg150 []float64
	dist150          []float64
	dense            []float64

	dataMinMax      []int
	subDataMinMax   []int
	avg30MinMax     []int
	avg150MinMax    []int
	avgMiddleMinMax []int

	resetMinMax []int

	cleanPosMinMax []int

	DIFF []float64
	DEA  []float64
	MACD []float64
	BAR  []float64
}

func (stock *Stock) resetData() {
	// 构造数据
	stock.dataOpen = make([]float64, 0)
	stock.dataClose = make([]float64, 0)

	//for i := 1; i <= 5; i++ {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 5; i > 2; i-- {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 2; i > 5; i++ {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	//for i := 5; i > 1; i-- {
	//	stock.dataOpen = append(stock.dataOpen, float64(i+1))
	//	stock.dataClose = append(stock.dataClose, float64(i))
	//}
	for i := 1; i < 10; i++ {
		stock.dataOpen = append(stock.dataOpen, float64(i))
		stock.dataClose = append(stock.dataClose, float64(i+1))
	}
	for i := 10; i > 5; i-- {
		stock.dataOpen = append(stock.dataOpen, float64(i))
		stock.dataClose = append(stock.dataClose, float64(i-1))
	}
	for i := 5; i < 15; i++ {
		stock.dataOpen = append(stock.dataOpen, float64(i))
		stock.dataClose = append(stock.dataClose, float64(i+1))
	}
	for i := 15; i > 3; i-- {
		stock.dataOpen = append(stock.dataOpen, float64(i))
		stock.dataClose = append(stock.dataClose, float64(i-1))
	}
	for i := 3; i < 20; i++ {
		stock.dataOpen = append(stock.dataOpen, float64(i))
		stock.dataClose = append(stock.dataClose, float64(i+1))
	}
	for i := 10; i < 3; i-- {
		stock.dataOpen = append(stock.dataOpen, float64(i))
		stock.dataClose = append(stock.dataClose, float64(i-1))
	}
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
		//if i > 300 {
		//	break
		//}

		valueOpen, err := strconv.ParseFloat(elems[1], 32)
		valueHigh, err := strconv.ParseFloat(elems[2], 32)
		valueClose, err := strconv.ParseFloat(elems[3], 32)
		valueLow, err := strconv.ParseFloat(elems[4], 32)
		if err != nil {
			log.Error("ParseFloat elems[3]:%s err:%v", elems[3], err)
		}
		stock.dataOpen = append(stock.dataOpen, valueOpen)
		stock.dataHigh = append(stock.dataHigh, valueHigh)
		stock.dataClose = append(stock.dataClose, valueClose)
		stock.dataLow = append(stock.dataLow, valueLow)
	}
	//
	//stock.resetData()

	log.Infof("filename: %s", filename)
	log.Infof("LoadAllData() stock.dataClose size: %d", len(stock.dataClose))

	//计算平均值
	for i := 0; i < len(stock.dataClose); i++ {
		val := get_pre_avg(stock.dataClose, i, 150)
		stock.avg150 = append(stock.avg150, val)
		val = get_pre_avg(stock.dataClose, i, 10)
		stock.avg10 = append(stock.avg10, val)
		val = get_pre_avg(stock.dataClose, i, 30)
		stock.avg30 = append(stock.avg30, val)

		//初始化
		stock.dataMinMax = append(stock.dataMinMax, 0)
		stock.subDataMinMax = append(stock.dataMinMax, 0)
		//stock.normalizedAvg150 = append(stock.normalizedAvg150, 0)
		//stock.dist150 = append(stock.dist150, 0)
		//stock.dense = append(stock.dense, 0)
	}
}

func copyStock(stock *Stock, start int, end int) *Stock {
	st := &Stock{
		dataOpen:      stock.dataOpen[start:end],
		dataClose:     stock.dataClose[start:end],
		dataHigh:      stock.dataHigh[start:end],
		dataLow:       stock.dataLow[start:end],
		avg10:         stock.avg10[start:end],
		dataMinMax:    stock.dataMinMax[start:end],
		subDataMinMax: stock.subDataMinMax[start:end],
	}
	return st
}
