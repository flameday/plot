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

type Rect struct {
	left      float64
	top       float64
	right     float64
	bottom    float64
	leftFlag  int
	rightFlag int
}

type avgContext struct {
	State          string
	Action         string
	Sell_stop      Rect
	Buy_stop       Rect
	High_Low_Min   float64
	Low_High_Max   float64
	Sell_Min_Value float64
	Buy_Max_Value  float64
	profit         float64
	buy            float64
	sell           float64
}

func (ac *avgContext) Show() string {
	s := "State : " + ac.State + " "
	s += "Action: " + ac.Action + " "
	if ac.Action == ACTION_BUY {
		s += fmt.Sprintf(" (%d, %.2f)->(%d, %.2f)",
			int(ac.Buy_stop.left),
			ac.Buy_stop.top,
			int(ac.Buy_stop.right),
			ac.Buy_stop.bottom)
	} else if ac.Action == ACTION_SELL {
		s += fmt.Sprintf(" (%d, %.2f)->(%d, %.2f)",
			int(ac.Sell_stop.left),
			ac.Sell_stop.top,
			int(ac.Sell_stop.right),
			ac.Sell_stop.bottom)
	} else {
		s += " Invalid Stop"
	}
	return s
}

//func (stock *Stock) GetDist() {
//	//计算距离 dist150
//	//均一化
//
//	filename := fmt.Sprintf("/Users/xinmei365/stock/dist.png")
//	p, _ := plot.New()
//	t := time.Now()
//	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
//	p.X.Label.Text = "GetDist"
//	p.Y.Label.Text = "dist"
//
//	drawData(p, stock.dataClose[0:400], 1, purple)
//	drawData(p, stock.avg150[0:400], 1, black)
//
//	//for i := 0; i < len(stock.dataClose); i++ {
//	for i := 0; i < 400; i++ {
//		start := i - 300
//		end := i + 1
//		if start < 0 {
//			start = 0
//			end = start + 300 + 1
//		}
//		min1, max1 := get_min_max(stock.dataClose[start:end])
//		min2, max2 := get_min_max(stock.avg150[start:end])
//		minY := math.Min(min1, min2)
//		maxY := math.Max(max1, max2)
//		deltaY := 50 / (maxY - minY + 0.01)
//
//		y := stock.dataClose[i]
//		normalizedY := (y - minY) * deltaY
//
//		// 找最小值
//		minDist := 10000000000.0
//		minPos := 10000000000
//		for pos := i - 60; pos < i+60; pos++ {
//			if pos < 0 || pos >= len(stock.dataClose) {
//				continue
//			}
//			//计算 normalizedAvg150
//			stock.normalizedAvg150[pos] = (stock.avg150[pos] - minY) * deltaY
//			//原地更新距离
//			stock.normalizedAvg150[pos] = float64((pos - i) * (pos - i))
//			stock.normalizedAvg150[pos] += (stock.normalizedAvg150[pos] - normalizedY) * (stock.normalizedAvg150[pos] - normalizedY)
//			if minDist > stock.normalizedAvg150[pos] {
//				minDist = stock.normalizedAvg150[pos]
//				minPos = pos
//			}
//		}
//		// 最小值作为距离
//		stock.dist150[i] = minDist
//
//		//drawLine(p, float64(i), stock.dataClose[i], float64(minPos), stock.normalizedAvg150[minPos])
//		drawLine(p, float64(i), stock.dataClose[i], float64(minPos), stock.avg150[minPos])
//
//		////取正负号
//		//if stock.dataClose[i] < stock.normalizedAvg150[i] {
//		//	stock.dist150[i] = -1 * stock.dist150[i]
//		//}
//		//// 遍历求最大小值
//		//tmp := make([]int, 0)
//		//caculateMinMax(stock.normalizedAvg150[i-60:i+60], &tmp, 12)
//		//for
//
//	}
//	p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
//}
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
func getSubWave(stock *Stock, pre int, post int) {

}
func isLow(data []float64, index int, length int) bool {
	//前后端点不计算大小值
	if index <= length/2 || index >= len(data)-length/2 {
		return false
	}

	for i := 1; i <= length/2; i++ {
		if data[index-i] < data[index] {
			return false
		}
		if data[index+i] < data[index] {
			return false
		}
	}
	return true
}
func isHigh(data []float64, index int, length int) bool {
	//前后端点不计算大小值
	if index <= length/2 || index >= len(data)-length/2 {
		return false
	}

	for i := 1; i <= length/2; i++ {
		if data[index-i] > data[index] {
			return false
		}
		if data[index+i] > data[index] {
			return false
		}
	}
	return true
}
func isRising(stock *Stock, pre int, post int) bool {
	if stock.dataMinMax[pre] < 0 && stock.dataMinMax[post] > 0 {
		return true
	}
	return false
}
func isFall(stock *Stock, pre int, post int) bool {
	if stock.dataMinMax[pre] > 0 && stock.dataMinMax[post] < 0 {
		return true
	}
	return false
}
func getSmallWave(stock *Stock, pre int, post int) {
	if isRising(stock, pre, post) {
		//先求高点
		for i := pre + 1; i < post; i++ {
			if isHigh(stock.dataHigh, i, 5) {
				stock.subDataMinMax[i] = 2
			}
		}
		//遍历高点，求后面的低点
		for i := pre + 1; i < post; i++ {
			if stock.subDataMinMax[i] == 2 {
				next := findNextIndex(stock.subDataMinMax, i+1, 2)
				if next == -1 {
					next = post
				}
				for k := i + 2; k < next; k++ {
					if isLow(stock.dataLow, k, 5) {
						//判断高低
						if stock.dataHigh[i] > stock.dataHigh[k] && stock.dataLow[i] > stock.dataLow[k] {
							stock.subDataMinMax[k] = -2
						}
					}
				}
			}
		}
	} else if isFall(stock, pre, post) {
		for i := pre + 1; i < post; i++ {
			if isLow(stock.dataLow, i, 5) {
				stock.subDataMinMax[i] = -2
			}
		}
		//遍历低点，求后面的高点
		for i := pre + 1; i < post; i++ {
			if stock.subDataMinMax[i] == -2 {
				next := findNextIndex(stock.subDataMinMax, i+1, -2)
				if next == -1 {
					next = post
				}
				for k := i + 2; k < next; k++ {
					if isHigh(stock.dataLow, k, 5) {
						//判断高低
						if stock.dataHigh[i] < stock.dataHigh[k] && stock.dataLow[i] < stock.dataLow[k] {
							stock.subDataMinMax[k] = 2
						}
					}
				}
			}
		}
	}
	//合并相邻的高点、合并相邻的低点
	lastPos := -1
	for i := pre + 1; i < post; i++ {
		if (stock.subDataMinMax[i] == 2) || (stock.subDataMinMax[i] == -2) {
			if lastPos == -1 || stock.subDataMinMax[lastPos] != stock.subDataMinMax[i] {
				lastPos = i
			} else if stock.subDataMinMax[lastPos] == stock.subDataMinMax[i] {
				if stock.subDataMinMax[lastPos] == 2 {
					if stock.dataHigh[lastPos] <= stock.dataHigh[i] {
						stock.subDataMinMax[lastPos] = 0

						lastPos = i
					} else {
						stock.subDataMinMax[i] = 0
						//keep lastPos
					}
				} else if stock.subDataMinMax[lastPos] == -2 {
					if stock.dataLow[lastPos] >= stock.dataLow[i] {
						stock.subDataMinMax[lastPos] = 0

						lastPos = i
					} else {
						stock.subDataMinMax[i] = 0
					}
				}
			}
		}
	}
}
