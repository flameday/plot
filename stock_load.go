package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"strconv"
	"strings"
)

type Stock struct {
	data      []float64
	avg6      []float64
	avg30     []float64
	avg150    []float64
	avgMiddle []float64

	dataMinMax      []int
	avg30MinMax     []int
	avg150MinMax    []int
	avgMiddleMinMax []int

	resetMinMax []int

	cleanPosMinMax []int
}

func (stock *Stock) LoadData() []float64 {
	// 读文本数据
	b, err := ioutil.ReadFile("/Users/xinmei365/stock_data_history/day/data/000002.csv")
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Println(b)
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

		// 解析数值
		val, err := strconv.ParseFloat(elems[3], 32)
		if err != nil {
			log.Error("ParseFloat elems[3]:%s err:%v", elems[3], err)
		}
		//
		if i > 2500 && i < 3500 {
			//if i < 500 {
			stock.data = append(stock.data, val)
		}
	}
	for i := 0; i < 10; i++ {
		log.Infof("[%d] %f", i, stock.data[i])
	}

	// 构造数据
	//stock.data = make([]float64, 0)
	//
	//for i := 1; i <= 20; i++ {
	//	stock.data = append(stock.data, float64(i))
	//}
	//for i := 20; i > 1; i-- {
	//	stock.data = append(stock.data, float64(i))
	//}
	//for i := 1; i < 20; i++ {
	//	stock.data = append(stock.data, float64(i))
	//}
	//for i := 20; i > 1; i-- {
	//	stock.data = append(stock.data, float64(i))
	//}
	//for i := 1; i < 20; i++ {
	//	stock.data = append(stock.data, float64(i))
	//}
	//for i := 20; i > 3; i-- {
	//	stock.data = append(stock.data, float64(i))
	//}
	//for i := 3; i < 20; i++ {
	//	stock.data = append(stock.data, float64(i))
	//}
	//for i := 10; i < 3; i-- {
	//	stock.data = append(stock.data, float64(i))
	//}

	//计算平均值
	for i := 0; i < len(stock.data); i++ {
		val := get_pre_avg(stock.data, i, 6)
		stock.avg6 = append(stock.avg6, val)

		val = get_pre_avg(stock.data, i, 30)
		stock.avg30 = append(stock.avg30, val)

		val = get_pre_avg(stock.data, i, 150)
		stock.avg150 = append(stock.avg150, val)

		val = get_middle_avg(stock.data, i, 30)
		stock.avgMiddle = append(stock.avgMiddle, val)
	}
	//
	log.Infof("data size:%d", len(stock.data))
	log.Infof("avg size:%d", len(stock.avg150))
	//局部最大值
	caculateMinMax(stock.data, &stock.dataMinMax, 30)
	caculateMinMax(stock.avg30, &stock.avg30MinMax, 30)
	caculateMinMax(stock.avg150, &stock.avg150MinMax, 150)
	caculateMinMax(stock.avgMiddle, &stock.avgMiddleMinMax, 30)
	//先使用平均值的minMax，后调整
	caculateMinMax(stock.avgMiddle, &stock.resetMinMax, 30)
	locateMax(stock.data, stock.resetMinMax, 61)
	locateMin(stock.data, stock.resetMinMax, 61)
	filter_max(stock.data, stock.resetMinMax)
	filter_min(stock.data, stock.resetMinMax)

	return stock.data
}

func caculateMinMax(data []float64, minMax *[]int, length int) {
	for i := 0; i < len(data); i++ {
		flag := 0
		if isMax(data, i, length) {
			flag = 1
		} else if isMin(data, i, length) {
			flag = -1
		}
		*minMax = append(*minMax, flag)
	}
}
