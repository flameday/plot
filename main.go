package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"image/color"
	"runtime/debug"
	"time"
)

//简单的逻辑
const (
	STATE_WAIT_FLAG = 0
	STATE_BUY_FLAG  = 1
	STATE_SELL_FLAG = 2
)

var (
	white    color.Color = color.RGBA{255, 255, 255, 255}
	blue     color.Color = color.RGBA{0, 0, 255, 255}
	red      color.Color = color.RGBA{255, 0, 0, 255}
	dark_red color.Color = color.RGBA{139, 0, 0, 255}
	green    color.Color = color.RGBA{0, 255, 0, 255}
	pink     color.Color = color.RGBA{255, 192, 203, 255}
	orange   color.Color = color.RGBA{255, 165, 0, 255}
	black    color.Color = color.RGBA{0, 0, 0, 255}
	gold     color.Color = color.RGBA{255, 215, 0, 255}
	yellow   color.Color = color.RGBA{255, 255, 0, 255}
	purple   color.Color = color.RGBA{128, 0, 128, 255}
	magenta  color.Color = color.RGBA{255, 0, 255, 255}
	olive    color.Color = color.RGBA{128, 128, 0, 255}
	gray     color.Color = color.RGBA{172, 172, 172, 255}

	picwidth       float64 = 512 * 2
	picheight      float64 = 384 * 2
	MAX_VALUE_FLAG         = 1
	MIN_VALUE_FLAG         = -1

	buy_stop  float64
	sell_stop float64
)

func work(stock *Stock, left int, right int) {
	ok, left, right := stock.LoadData(left, right)
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
	drawData(p, stock.avg30, 1, blue)
	drawData(p, stock.avg150, 3, green)

	//drawMinMax(p, stock.avgMiddle, stock.avgMiddleMinMax, 1, 3, blue)
	//drawMinMax(p, stock.avgMiddle, stock.avgMiddleMinMax, -1, 2, purple)
	drawMinMax(p, stock.dataClose, stock.resetMinMax, 1, 2, gray)
	drawMinMax(p, stock.dataClose, stock.resetMinMax, -1, 2, gray)

	aimArr := make([]int, 0)
	aimArr = append(aimArr, 1)
	aimArr = append(aimArr, -1)
	drawMinMax2(p, stock.dataClose, stock.resetMinMax, aimArr, 2, black)

	//drawMinMax(p, stock.dataClose, stock.flagArea, -1, 3, green)
	drawData(p, stock.relateCntArray, 2, green)

	name := fmt.Sprintf("/Users/xinmei365/stock/price_%d_%d.png", left, right)
	if right >= len(stock.dataClose) {
		name = fmt.Sprintf("/Users/xinmei365/stock/price_all.png")
	}

	p.Save(vg.Length(picwidth), vg.Length(picheight), name)

	//stock.LoadData("/Users/xinmei365/stock_data_history/day/dataClose/000002.csv")
	//http.HandleFunc("/", RrawPicture)
	//http.ListenAndServe(":999", nil)
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

	//// 绘图
	//for i := 0; i < 100; i++ {
	//	stock := Stock{}
	//	left := i*1000 - 200
	//	right := (i+1)*1000 + 200
	//	work(&stock, left, right)
	//}
	//// 最后加一个完整的图
	//stock := Stock{}
	//work(&stock, 0, 10000)
	// 遍历模拟
	action_state := STATE_WAIT_FLAG
	for i := 1; i < 10000; i++ {
		stock := Stock{}
		ok, _, _ := stock.LoadData(0, i)
		if !ok {
			return
		}
		changeAction(&action_state, stock.dataClose, stock.resetMinMax, i)
	}
}
