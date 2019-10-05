package main

import (
	"flag"
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
	gray           color.Color = color.RGBA{196, 196, 196, 255}
	colorArray                 = []color.Color{red, blue, black, yellow, orange, gold, purple, magenta, olive, gray}
	picwidth       float64     = 512 * 2
	picheight      float64     = 384 * 2
	MAX_VALUE_FLAG             = 1
	MIN_VALUE_FLAG             = -1

	buy_stop  float64
	sell_stop float64

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

	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "/Users/xinmei365/node_modules/editor.md/fonts/editormd-logo.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 125, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
	text     = string("JOJO hoho")
)

func run(ac *avgContext, p *plot.Plot, stock *Stock, filename string, pos int) bool {

	//_, st := getAllRect(stock)
	curPos := len(stock.dataClose) - 1

	if ac.State == STATE_UNKOWN {
		ok, revert, change := isValidInit(ac, stock)
		if ok && revert && change {
			//log.Infof("ac: %s", ac.Show())

			drawPoint(p, float64(curPos), stock.dataClose[curPos], 20, red)
			if ac.Action == ACTION_BUY {
				drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
			} else if ac.Action == ACTION_SELL {
				drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
			}

			p.X.Label.Text = ac.State + " " + ac.Action
			//p.Save(vg.Length(picwidth), vg.Length(picheight), filename)

			return true
		}
	} else if ac.State != STATE_UNKOWN {
		ok, revert, change, arr := forwardState(ac, stock)
		for _, r := range arr {
			drawRectangle(p, r.left, r.top, r.right, r.bottom, olive)
			if r.leftFlag == -1 {
				drawPoint(p, r.left, r.top, 10, red)
			}
			if r.leftFlag == 1 {
				drawPoint(p, r.left, r.top, 15, blue)
			}
		}
		//log.Infof("pos:%d ok:%v", pos, ok)
		if ok && (revert || change) {

			p.X.Label.Text = ac.State + " " + ac.Action
			drawPoint(p, float64(curPos), stock.dataClose[curPos], 20, black)
			if ac.Action == ACTION_BUY {
				drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
			} else if ac.Action == ACTION_SELL {
				drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
			}
			p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
		}
	} else {
		log.Infof("ignore:%d", pos)
	}
	return false
}
func main() {
	defer func() {
		if err := recover(); err != nil {

			log.Infof("err:%v", err)
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
	for index := 0; index < len(dstArray); index += 5 {
		if index < 15 {
			continue
		}
		if index > 15 {
			break
		}

		textfile := dstArray[index]
		stockBig := Stock{}
		stockBig.LoadAllData(textfile)

		//ac := &avgContext{
		//	State:          STATE_UNKOWN,
		//	profit:         0.0,
		//	Sell_Min_Value: -1,
		//	Buy_Max_Value:  -1,
		//}

		tmpArr := make([]float64, 0)
		//for i := 0; i < 500; i++ {
		//	getWave(&stockBig, i)
		//}
		for i := 1; i < 301; i += 1 {
			//for i := 1; i < 100; i += 1 {
			p, _ := plot.New()
			t := time.Now()
			p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
			p.X.Label.Text = "drawWithRect"
			p.Y.Label.Text = "Price"

			start := i - 300
			end := i + 1
			if start < 0 {
				start = 0
				end = start + 300 + 1
			}

			////1， 绘制底图
			//getAllRect(&stockBig)
			//drawData(p, stockBig.dataHigh[start:end], 2, red)
			//drawData(p, stockBig.dataLow[start:end], 1, gray)
			for k := start; k < end; k++ {
				drawLine(p, float64(k-start), stockBig.dataLow[k], float64(k-start), stockBig.dataHigh[k])
			}
			//drawData(p, stockBig.dataLow[start:end], 2, yellow)
			drawData(p, stockBig.avg10[start:end], 1, purple)

			//st := copyStock(&stockBig, start, i+1)
			st := copyStock(&stockBig, start, end)
			arr, _ := getAllRect(st)
			if len(arr) > 0 {
				//	for _, r := range arr {
				//		drawRectangle(p, r.left, r.top, r.right, r.bottom, gray)
				//	}
			}
			drawAllSubMinMax(p, st, 2, blue)
			drawAllMinMax(p, st, 2, black)

			filename := fmt.Sprintf("/Users/xinmei365/stock/%03d_%03d.png", index, i)
			//ret := run(ac, p, st, filename, i)
			//if ret {
			//	tmpArr = append(tmpArr, ac.profit)
			//	log.Infof("[%d] profit:%f", i, ac.profit)
			//}
			//if i%300 == 0 {
			p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
			//}
			break
		}
		//
		drawPic(tmpArr, "Count", "Profit", fmt.Sprintf("/Users/xinmei365/profilt_%d.png", index))
	}
}

func drawPic(data []float64, xlabel string, ylabel string, filename string) {
	p, _ := plot.New()
	t := time.Now()
	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	drawData(p, data, 1, black)

	p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
}
