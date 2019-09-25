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

	STATE_UNKOWN   = 0
	STATE_INIT     = 1
	STATE_NEW_HIGH = 2
	STATE_NEW_LOW  = 4

	ACTION_NONE = "ACTION_NONE"
	ACTION_BUY  = "ACTION_BUY"
	ACTION_SELL = "ACTION_SELL"
)

func getState(state int) string {
	if state == 0 {
		return "STATE_UNKOWN"
	}
	if state == 1 {
		return "STATE_INIT"
	}
	if state == 2 {
		return "STATE_NEW_HIGH"
	}
	if state == 4 {
		return "STATE_NEW_LOW"
	}
	log.Errorf("getState error, state: %d", state)
	return ""
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
		if index > 12 {
			break
		}

		filename := dstArray[index]
		for i := 0; i < 10; i++ {
			stock := Stock{}
			//left := i*500 - 100
			//right := (i + 1) * 500
			//work(filename, &stock, index, left, right)
			stock.LoadAllData(filename)

			// draw
			// 创建 plog
			p, _ := plot.New()
			t := time.Now()

			p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
			p.X.Label.Text = "drawWithRect"
			p.Y.Label.Text = "Price"

			drawData(p, stock.dataClose[0:500], 1, red)

			ac := &avgContext{
				State: STATE_UNKOWN,
			}
			for pos := 10; pos < 500; pos += 10 {
				arr, st := getAllRect(stock.dataClose[0:pos])
				for _, r := range arr {
					drawRectangle(p, r.left, r.top, r.right, r.bottom)
					//drawLine(p, r.left, r.top, r.right, r.bottom)
					//drawLine(p, r.left, r.bottom, r.right, r.top)
				}
				// 高低点
				// 保存图片
				filename := fmt.Sprintf("/Users/xinmei365/stock/%d_%d_%d.png", index, i, pos)

				drawMinMax(p, st.dataClose[0:pos], st.dataMinMax[0:pos], 1, 3, gray)

				if ac.State == STATE_UNKOWN {
					ok := isValidInit(ac, st.dataClose[0:pos])
					if ok {
						log.Infof("state: %s ac: %v", *ac)
						drawPoint2(p, float64(pos), st.dataClose[pos-1], 12, red)
						//p.Save(vg.Length(picwidth), vg.Length(picheight), filename)

						continue
					}
				}

				if ac.State != STATE_UNKOWN {
					ok := forwardState(ac, st.dataClose[0:pos])
					if ok {
						clr := red
						if ac.State == STATE_NEW_HIGH+STATE_NEW_LOW {
							clr = purple
							drawPoint2(p, float64(pos), st.dataClose[pos-1], 20, clr)
							p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
						}
					}
				}
			}

			break
		}
	}
}
