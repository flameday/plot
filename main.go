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

func main() {
	defer func() {
		if err := recover(); err != nil {

			defer log.Flush()

			debug.PrintStack()
		}
	}()

	// log
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

			// draw
			// 创建 plog
			p, _ := plot.New()
			t := time.Now()

			p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
			p.X.Label.Text = "drawWithRect"
			p.Y.Label.Text = "Price"

			drawData(p, stock.dataClose[0:500], 1, red)
			for _, r := range arr {
				drawRectangle(p, r.left, r.top, r.right, r.bottom)
				//drawLine(p, r.left, r.top, r.right, r.bottom)
				//drawLine(p, r.left, r.bottom, r.right, r.top)

			}
			// 高低点
			drawMinMax(p, st.dataClose[0:500], st.dataMinMax[0:500], 1, 3, gray)

			// 保存图片
			p.Save(vg.Length(picwidth), vg.Length(picheight), fmt.Sprintf("/Users/xinmei365/stock/%d_%d.png", index, i))

			break
		}
	}
}
