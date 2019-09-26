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

	STATE_UNKOWN              = "STATE_UNKOWN"
	STATE_NEW_HIGH            = "STATE_NEW_HIGH"
	STATE_NEW_LOW             = "STATE_NEW_LOW"
	STATE_NEW_LOW__NEW_HIGH_0 = "STATE_NEW_LOW__NEW_HIGH_0"
	STATE_NEW_LOW__NEW_HIGH_1 = "STATE_NEW_LOW__NEW_HIGH_1"
	STATE_NEW_HIGH__NEW_LOW_0 = "STATE_NEW_HIGH__NEW_LOW_0"
	STATE_NEW_HIGH__NEW_LOW_1 = "STATE_NEW_HIGH__NEW_LOW_1"

	ACTION_NONE = "ACTION_NONE"
	ACTION_BUY  = "ACTION_BUY"
	ACTION_SELL = "ACTION_SELL"
)

func run(p *plot.Plot, data []float64, filename string) {

	ac := &avgContext{
		State: STATE_UNKOWN,
	}

	_, st := getAllRect(data)
	curPos := len(data) - 1

	//log.Infof("len(arr):%v", len(arr))
	//for _, r := range arr {
	//	//drawRectangle(p, r.left, r.top, r.right, r.bottom, gray)
	//	//drawLine(p, r.left, r.top, r.right, r.bottom)
	//	//drawLine(p, r.left, r.bottom, r.right, r.top)
	//}
	//drawMinMax(p, st.dataClose[0:pos], st.dataMinMax[0:pos], 1, 1, yellow)
	p.Save(vg.Length(picwidth), vg.Length(picheight), filename)

	if ac.State == STATE_UNKOWN {
		ok, revert, change := isValidInit(ac, st.dataClose)
		if ok && revert && change {
			log.Infof("ac: %s", ac.Show())

			drawPoint2(p, float64(curPos), st.dataClose[curPos], 20, red)
			if ac.Action == ACTION_BUY {
				drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
			} else if ac.Action == ACTION_SELL {
				drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
			}

			p.X.Label.Text = ac.State + " " + ac.Action
			p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
		}
	} else if ac.State != STATE_UNKOWN {
		ok, revert, change := forwardState(ac, st.dataClose)
		if ok {
			p.X.Label.Text = ac.State + " " + ac.Action
			if revert {
				log.Infof("ac: %s", ac.Show())

				drawPoint2(p, float64(curPos), st.dataClose[curPos], 20, black)

				if ac.Action == ACTION_BUY {
					drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
				} else if ac.Action == ACTION_SELL {
					drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
				}
				p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
			} else if change {
				log.Infof("ac: %s", ac.Show())

				drawPoint2(p, float64(curPos), st.dataClose[curPos], 20, black)

				if ac.Action == ACTION_BUY {
					drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
				} else if ac.Action == ACTION_SELL {
					drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
				}
				p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
			}
		}
	}
}
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
		if index > 10 {
			break
		}

		filename := dstArray[index]
		stock := Stock{}
		stock.LoadAllData(filename)

		for i := 1; i < len(stock.dataClose); i += 30 {
			p, _ := plot.New()
			t := time.Now()
			p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
			p.X.Label.Text = "drawWithRect"
			p.Y.Label.Text = "Price"

			start := 0
			end := 300
			//pos := 1
			if i >= 300 {
				start = i - 300 + 1
				end = i + 1
				//pos = 299
			}

			//1， 绘制底图
			drawData(p, stock.dataClose[start:end], 1, pink)
			filename := fmt.Sprintf("/Users/xinmei365/stock/%03d_%03d.png", index, i)
			run(p, stock.dataClose[start:i+1], filename)
			//p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
		}
	}
}
