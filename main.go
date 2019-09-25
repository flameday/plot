package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"image/color"
	"math"
	"runtime/debug"
	"time"
)

//简单的逻辑
const ()

var (
	white          color.Color = color.RGBA{255, 255, 255, 255}
	blue           color.Color = color.RGBA{0, 0, 255, 255}
	red            color.Color = color.RGBA{255, 0, 0, 255}
	dark_red       color.Color = color.RGBA{139, 0, 0, 255}
	green          color.Color = color.RGBA{0, 255, 0, 255}
	pink           color.Color = color.RGBA{255, 192, 203, 255}
	orange         color.Color = color.RGBA{255, 165, 0, 255}
	black          color.Color = color.RGBA{0, 0, 0, 255}
	gold           color.Color = color.RGBA{255, 215, 0, 255}
	yellow         color.Color = color.RGBA{255, 255, 0, 255}
	purple         color.Color = color.RGBA{128, 0, 128, 255}
	magenta        color.Color = color.RGBA{255, 0, 255, 255}
	olive          color.Color = color.RGBA{128, 128, 0, 255}
	gray           color.Color = color.RGBA{172, 172, 172, 255}
	colorArray                 = []color.Color{red, blue, black, yellow, orange, gold, purple, magenta, olive, gray}
	picwidth       float64     = 512 * 2
	picheight      float64     = 384 * 2
	MAX_VALUE_FLAG             = 1
	MIN_VALUE_FLAG             = -1

	buy_stop  float64
	sell_stop float64
	flag      int = 0
)

func getDense(data []float64, start int, end int) float64 {
	diff := math.Abs(data[end] - data[start])
	dense := diff / float64(end-start)
	return dense
}
func getDistanceArray(data []float64, startPos int, aimPos int) []float64 {
	arr := make([]float64, 0)
	if startPos < aimPos {
		for i := aimPos - 1; i >= startPos; i-- {
			val := data[i] - data[aimPos]
			arr = append(arr, val*val)
		}
	}
	if startPos > aimPos {
		for i := aimPos + 1; i <= startPos; i++ {
			val := data[i] - data[aimPos]
			arr = append(arr, val*val)
		}
	}
	return arr
}

func compareDistanceArray(leftArr []float64, rightArr []float64) int {
	cnt := 0
	for i := 0; i < 8; i++ {
		if i >= len(leftArr) {
			break
		}
		if i >= len(rightArr) {
			break
		}
		if leftArr[i] < rightArr[i] {
			cnt += 1
		} else {
			cnt -= 1
		}
	}
	return cnt
}
func getRectangle(data []float64, posLeft int, posMiddle int, posRight int) (int, int) {
	leftArr := getDistanceArray(data, posLeft, posMiddle)
	rightArr := getDistanceArray(data, posMiddle, posRight)
	val := compareDistanceArray(leftArr, rightArr)
	if val > 0 {
		return posLeft, posMiddle
	}
	return posMiddle, posRight

	//var01 := getDense(data, posLeft, posMiddle)
	//var02 := getDense(data, posMiddle, posRight)
	//log.Infof("var01, var02:%f %f", var01, var02)
	//if var01 <= var02 {
	//	log.Infof("posLeft, posMiddle:[%d, %d]", posLeft, posMiddle)
	//	return posLeft, posMiddle
	//}
	//log.Infof("posMiddle, posRight:[%d, %d]", posMiddle, posRight)
	//return posMiddle, posRight
}
func getRectangle3(data []float64, posLeft int, posMiddle int, posRight int) (int, int) {
	var01 := getDense(data, posLeft, posMiddle)
	var02 := getDense(data, posMiddle, posRight)
	log.Infof("var01, var02:%f %f", var01, var02)
	if var01 <= var02 {
		log.Infof("posLeft, posMiddle:[%d, %d]", posLeft, posMiddle)
		return posLeft, posMiddle
	}
	log.Infof("posMiddle, posRight:[%d, %d]", posMiddle, posRight)
	return posMiddle, posRight
}

func getRectangle2(data []float64, posLeft int, posMiddle int, posRight int) (int, int) {
	var01 := getVariance(data, posLeft, posMiddle+1)
	var02 := getVariance(data, posMiddle, posRight+1)
	if var01 >= var02 {
		return posLeft, posMiddle
	}
	return posMiddle, posRight
}

func drawWithRect(data []float64, arr []Rect) *plot.Plot {
	//创建 plog
	p, _ := plot.New()
	t := time.Now()

	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
	p.X.Label.Text = "drawWithRect"
	p.Y.Label.Text = "Price"

	drawData(p, data, 1, red)
	for _, r := range arr {
		drawRectangle(p, r.left, r.top, r.right, r.bottom)
		//drawLine(p, r.left, r.top, r.right, r.bottom)
		//drawLine(p, r.left, r.bottom, r.right, r.top)

	}

	return p
}

func work(filename string, stock *Stock, index int, left int, right int) {
	flagOver := false
	if right >= 10000 {
		flagOver = true
	}

	ok, left, right := stock.LoadData(filename, left, right)
	if !ok {
		return
	}

	//创建 plog
	p, _ := plot.New()
	t := time.Now()

	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
	p.X.Label.Text = "Quantity Demand"
	p.Y.Label.Text = "Price"

	drawData(p, stock.dataClose, 3, red)
	//drawData(p, stock.dataOpen, 2, red)
	//drawData(p, stock.avgMiddle, 2, dark_red)
	//drawData(p, stock.avg6, 1, green)
	//drawData(p, stock.avg30, 1, blue)
	//drawData(p, stock.avg150, 3, green)

	//drawMinMax(p, stock.dataClose, stock.dataMinMax, 1, 3, blue)
	//drawMinMax(p, stock.dataClose, stock.dataMinMax, -1, 3, blue)
	//drawMinMax(p, stock.avgMiddle, stock.avgMiddleMinMax, 1, 3, blue)
	//drawMinMax(p, stock.avgMiddle, stock.avgMiddleMinMax, -1, 2, purple)
	//drawMinMax(p, stock.dataClose, stock.resetMinMax, 1, 2, gray)
	//drawMinMax(p, stock.dataClose, stock.resetMinMax, -1, 2, gray)

	//rectangle
	for i := 1; i < len(stock.dataMinMax)-1; i++ {
		if stock.dataMinMax[i] == 1 {
			preMin := findPreIndex(stock.dataMinMax, i, -1)
			postMin := findNextIndex(stock.dataMinMax, i, -1)
			if preMin != -1 && postMin != -1 {
				x1, x2 := getRectangle(stock.dataClose, preMin, i, postMin)
				//drawRectangle(p, float64(x1), stock.dataClose[x1], float64(x2), stock.dataClose[x2])
				drawLine(p, float64(x1), stock.dataClose[x1], float64(x2), stock.dataClose[x2])
			}
		}
		if stock.dataMinMax[i] == -1 {
			preMax := findPreIndex(stock.dataMinMax, i, 1)
			postMax := findNextIndex(stock.dataMinMax, i, 1)
			if preMax != -1 && postMax != -1 {
				x1, x2 := getRectangle(stock.dataClose, preMax, i, postMax)
				//drawRectangle(p, float64(x1), stock.dataClose[x1], float64(x2), stock.dataClose[x2])
				drawLine(p, float64(x1), stock.dataClose[x1], float64(x2), stock.dataClose[x2])
			}
		}
	}

	//for i := 0; i < len(stock.dataClose); {
	//	pos := findBump(stock.dataClose, i)
	//	if pos != -1 {
	//		i = pos
	//		drawPoint(p, float64(i), stock.dataClose[i], 3)
	//		log.Infof("findBump pos: %d", pos)
	//	} else {
	//		i++
	//	}
	//}

	//aimArr := make([]int, 0)
	//aimArr = append(aimArr, 1)
	//aimArr = append(aimArr, -1)
	//drawMinMax2(p, stock.dataClose, stock.resetMinMax, aimArr, 2, black)

	//drawMinMax(p, stock.dataClose, stock.flagArea, -1, 3, green)
	//drawData(p, stock.relateCntArray, 2, green)

	name := fmt.Sprintf("/Users/xinmei365/stock/price_%d_%d_%d.png", index, left, right)
	if flagOver {
		name = fmt.Sprintf("/Users/xinmei365/stock/price_%d_all.png", index)
	}

	p.Save(vg.Length(picwidth), vg.Length(picheight), name)

	//stock.LoadData("/Users/xinmei365/stock_data_history/day/dataClose/000002.csv")
	//http.HandleFunc("/", RrawPicture)
	//http.ListenAndServe(":999", nil)
}

func findBump(data []float64, pos int) int {
	if pos < 0 {
		return -1
	}
	if pos+3 >= len(data) {
		return -1
	}

	slope0 := data[pos+1] - data[pos+0]
	slope1 := data[pos+2] - data[pos+1]
	flag := slope1 - slope0

	preSlope := slope1
	for i := pos + 3; i < len(data); i++ {
		curSlope := data[i] - data[i-1]
		if (curSlope-preSlope)*flag < 0 {
			return i
		}
		preSlope = curSlope
	}
	return -1
}

// 大家可以查看这个网址看看这个image包的使用方法 http://golang.org/doc/articles/image_draw.html
func main() {
	defer func() {
		if err := recover(); err != nil {

			defer log.Flush()

			debug.PrintStack()
		}
	}()
	logger, err := log.LoggerFromConfigAsFile("/Users/xinmei365/go/src/plot/conf/log.xml")
	if err != nil {
		fmt.Printf("parse config.xml error")
		log.Errorf("parse config.xml error")
	}
	log.ReplaceLogger(logger)
	defer log.Flush()

	// 绘图
	fileArray := make([]string, 0)
	dstArray, err := GetAllFile("/Users/xinmei365/stock_data_history/day/data/", fileArray)
	//filename := "/Users/xinmei365/stock_data_history/day/data/000002.csv"
	for index := 0; index < len(dstArray); index++ {
		if index < 10 {
			continue
		}
		if index > 21 {
			break
		}

		filename := dstArray[index]
		for i := 0; i < 10; i++ {
			stock := Stock{}
			//left := i*500 - 100
			//right := (i + 1) * 500
			//work(filename, &stock, index, left, right)
			stock.LoadAllData(filename)
			arr, st := getAllRect(stock.dataClose[0:500])
			//draw
			p := drawWithRect(st.dataClose[0:500], arr)
			drawMinMax(p, st.dataClose[0:500], st.dataMinMax[0:500], 1, 3, gray)

			p.Save(vg.Length(picwidth), vg.Length(picheight), fmt.Sprintf("/Users/xinmei365/stock/%d_%d.png", index, i))

			break
		}
		//// 最后加一个完整的图
		stock := Stock{}
		work(filename, &stock, index, 0, 10000)
	}

	// 遍历模拟
	//action_state := ACTION_WAIT_FLAG
	//for i := 1; i < 10000; i++ {
	//	stock := Stock{}
	//	ok, _, _ := stock.LoadData(0, i)
	//	if !ok {
	//		return
	//	}
	//	changeAction(&action_state, stock.dataClose, stock.resetMinMax, i)
	//}
}
