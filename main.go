package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"image/color"
	"time"
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

	picwidth       float64 = 512 * 2
	picheight      float64 = 384 * 2
	stock          Stock
	MAX_VALUE_FLAG = 1
	MIN_VALUE_FLAG = -1
)

// 大家可以查看这个网址看看这个image包的使用方法 http://golang.org/doc/articles/image_draw.html
func main() {
	logger, err := log.LoggerFromConfigAsFile("/Users/xinmei365/go/src/plot/conf/log.xml")
	if err != nil {
		fmt.Printf("parse config.xml error")
		log.Errorf("parse config.xml error")
	}
	log.ReplaceLogger(logger)
	defer log.Flush()

	stock.LoadData()

	//创建 plog
	p, _ := plot.New()
	t := time.Now()

	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
	p.X.Label.Text = "Quantity Demand"
	p.Y.Label.Text = "Price"

	drawData(p, stock.dataClose, 2, red)
	//drawData(p, stock.avgMiddle, 2, dark_red)
	//drawData(p, stock.avg6, 1, green)
	drawData(p, stock.avg30, 1, purple)
	drawData(p, stock.avg150, 5, yellow)

	//drawMinMax(p, stock.avgMiddle, stock.avgMiddleMinMax, 1, 3, blue)
	//drawMinMax(p, stock.avgMiddle, stock.avgMiddleMinMax, -1, 2, purple)
	drawMinMax(p, stock.dataClose, stock.resetMinMax, 1, 2, black)
	drawMinMax(p, stock.dataClose, stock.resetMinMax, -1, 2, black)

	//drawMinMax(p, stock.dataClose, stock.flagArea, -1, 3, green)
	drawData(p, stock.relateCntArray, 2, green)

	p.Save(vg.Length(picwidth), vg.Length(picheight), "/Users/xinmei365/price.png")

	//stock.LoadData("/Users/xinmei365/stock_data_history/day/dataClose/000002.csv")
	//http.HandleFunc("/", RrawPicture)
	//http.ListenAndServe(":999", nil)
}
