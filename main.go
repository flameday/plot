package main

import (
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
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

type exampleThumbnailer struct {
	color.Color
}

// Thumbnail fulfills the plot.Thumbnailer interface.
func (et exampleThumbnailer) Thumbnail(c *draw.Canvas) {
	pts := []vg.Point{
		{c.Min.X, c.Min.Y},
		{c.Min.X, c.Max.Y},
		{c.Max.X, c.Max.Y},
		{c.Max.X, c.Min.Y},
	}
	poly := c.ClipPolygonY(pts)
	c.FillPolygon(et.Color, poly)

	pts = append(pts, vg.Point{X: c.Min.X, Y: c.Min.Y})
	outline := c.ClipLinesY(pts)
	c.StrokeLines(draw.LineStyle{
		Color: color.Black,
		Width: vg.Points(1),
	}, outline...)
}

func run(ac *avgContext, p *plot.Plot, stock *Stock, filename string, pos int) bool {

	//log.Infof("pos:%d", pos)
	//if pos == 58 {
	//	log.Infof("58")
	//}
	_, st := getAllRect(stock)
	curPos := len(stock.dataClose) - 1

	if ac.State == STATE_UNKOWN {
		ok, revert, change := isValidInit(ac, st)
		if ok && revert && change {
			//log.Infof("ac: %s", ac.Show())

			drawPoint(p, float64(curPos), st.dataClose[curPos], 20, red)
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
		ok, revert, change, arr := forwardState(ac, st)
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
			drawPoint(p, float64(curPos), st.dataClose[curPos], 20, black)
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

//func drawInflection(stock *Stock, xlabel string, ylabel string, filename string) {
//	p, _ := plot.New()
//	t := time.Now()
//	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
//	p.X.Label.Text = xlabel
//	p.Y.Label.Text = ylabel
//
//	preRiseLength := 0.0
//	preRiseCount := 0
//	preDownLength := 0.0
//	preDownCount := 0
//	trendRise := true
//	trendData := make([]float64, 0)
//	// 初始化
//	preDiff := stock.dataClose[0] - stock.dataOpen[0]
//	if preDiff >= 0 {
//		trendRise = true
//		preRiseLength = preDiff
//		preRiseCount = 1
//	} else {
//		trendRise = false
//		preDownLength = preDiff
//		preDownCount = 1
//	}
//
//	if trendRise {
//		trendData = append(trendData, 1.0)
//	} else {
//		trendData = append(trendData, -1.0)
//	}
//
//	// 遍历
//	for i := 1; i < len(stock.dataOpen); i++ {
//
//		if i%2 == 0 {
//			//drawLine(p, float64(i), -1, float64(i), 20)
//		}
//
//		diff := stock.dataClose[i] - stock.dataOpen[i]
//		log.Infof("[%d] trendRise:%v diff:%f", i, trendRise, diff)
//		log.Infof("[%d] rise:%f %d down:%f %d", i, preRiseLength, preRiseCount, preDownLength, preDownCount)
//
//		if diff == 0 {
//			continue
//		}
//
//		if trendRise == true {
//			if diff >= 0 {
//				preRiseLength += math.Abs(diff)
//				preRiseCount += 1
//			} else if diff < 0 {
//				// 转向
//				if math.Abs(diff)+preDownLength > preRiseLength/float64(preRiseCount) {
//					drawPoint(p, float64(i), stock.dataOpen[i], 10, red)
//				}
//				// 转势
//				if math.Abs(diff)+preDownLength > 0.5*preRiseLength {
//					trendRise = false
//
//					preRiseLength = 0
//					preRiseCount = 0
//				}
//
//				preDownLength += math.Abs(diff)
//				preDownCount += 1
//				log.Infof("[%d] rise:%f %d down:%f %d\n", i, preRiseLength, preRiseCount, preDownLength, preDownCount)
//			}
//		} else if trendRise == false {
//			if diff <= 0 {
//				preDownLength += math.Abs(diff)
//				preDownCount += 1
//			} else if diff > 0 {
//				// 转向
//				if math.Abs(diff)+preRiseLength > preDownLength/float64(preDownCount) {
//					drawPoint(p, float64(i), stock.dataOpen[i], 15, blue)
//				}
//				// 转势
//				if math.Abs(diff)+preRiseLength > 0.5*preDownLength {
//					trendRise = true
//
//					preDownLength = 0
//					preDownCount = 0
//				}
//
//				preRiseLength += math.Abs(diff)
//				preRiseCount += 1
//				log.Infof("[%d] rise:%f %d down:%f %d\n", i, preRiseLength, preRiseCount, preDownLength, preDownCount)
//			}
//		}
//		if trendRise {
//			trendData = append(trendData, 0.1)
//		} else {
//			trendData = append(trendData, 0)
//		}
//	}
//
//	drawData(p, stock.dataOpen, 2, dark_red)
//	drawData(p, stock.dataClose, 2, gray)
//	drawData(p, trendData, 2, magenta)
//	p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
//}

//func drawSubBar(stock *Stock, xlabel string, ylabel string, filename string) {
//	p, _ := plot.New()
//	t := time.Now()
//	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
//	p.X.Label.Text = xlabel
//	p.Y.Label.Text = ylabel
//
//	dataBar := make([]float64, 0)
//	flagBar := make([]int, 0)
//	//totalBar := make([]float64, 0)
//	closeBar := make([]int, 0)
//	for i := 0; i < len(stock.dataOpen); i++ {
//		if stock.dataOpen[i] <= stock.dataClose[i] { //rising
//			//down - up - down
//			diff01 := stock.dataOpen[i] - stock.dataLow[i]
//			flag01 := -1
//			diff02 := stock.dataHigh[i] - stock.dataLow[i]
//			flag02 := 1
//			diff03 := stock.dataHigh[i] - stock.dataClose[i]
//			flag03 := -1
//			//if diff01 > 0 {
//			dataBar = append(dataBar, diff01)
//			flagBar = append(flagBar, flag01)
//			//}
//			//if diff02 > 0 {
//			dataBar = append(dataBar, diff02)
//			flagBar = append(flagBar, flag02)
//			//}
//			//if diff03 > 0 {
//			dataBar = append(dataBar, diff03)
//			flagBar = append(flagBar, flag03)
//			//}
//
//			closeBar = append(closeBar, 1)
//		}
//
//		if stock.dataOpen[i] > stock.dataClose[i] { //rising
//			//down - up - down
//			diff01 := stock.dataHigh[i] - stock.dataOpen[i]
//			flag01 := 1
//			diff02 := stock.dataHigh[i] - stock.dataLow[i]
//			flag02 := -1
//			diff03 := stock.dataClose[i] - stock.dataLow[i]
//			flag03 := 1
//			//if diff01 > 0 {
//			dataBar = append(dataBar, diff01)
//			flagBar = append(flagBar, flag01)
//			//}
//			//if diff02 > 0 {
//			dataBar = append(dataBar, diff02)
//			flagBar = append(flagBar, flag02)
//			//}
//			//if diff03 > 0 {
//			dataBar = append(dataBar, diff03)
//			flagBar = append(flagBar, flag03)
//			//}
//
//			closeBar = append(closeBar, -1)
//		}
//	}
//	for i := 1; i < len(dataBar); i++ {
//		//totalBar = append(totalBar, dataBar[i]*float64(flagBar[i]))
//		//totalBar = append(totalBar, dataBar[i]-dataBar[i-1])
//	}
//	//drawData(p, totalBar, 2, blue)
//	//drawMinMax(p, dataBar, flagBar, 1, 1, blue)
//	//drawMinMax(p, dataBar, flagBar, -1, 1, red)
//	drawMinMax(p, stock.dataClose, closeBar, 1, 1, blue)
//	drawMinMax(p, stock.dataClose, closeBar, -1, 1, red)
//
//	p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
//}

func drawPic(data []float64, xlabel string, ylabel string, filename string) {
	p, _ := plot.New()
	t := time.Now()
	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	drawData(p, data, 1, black)

	p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
}

//func merge(dataHigh []float64, dataLow []float64) ([]float64, []float64) {
//	if len(dataHigh) <= 0 {
//		return dataHigh, dataLow
//	}
//	if len(dataHigh) != len(dataLow) {
//		log.Errorf("merge error")
//	}
//
//	newHigh := make([]float64, 0)
//	newLow := make([]float64, 0)
//	newHigh = append(newHigh, dataHigh[0])
//	newLow = append(newLow, dataLow[0])
//
//	for i := 1; i < len(dataHigh); i++ {
//		if (dataHigh[i-1] >= dataHigh[i]) && (dataLow[i-1] <= dataLow[i]) {
//			newHigh[len(newHigh)-1] = math.Max(newHigh[len(newHigh)-1], dataHigh[i])
//			newLow[len(newLow)-1] = math.Min(newLow[len(newLow)-1], dataLow[i])
//		} else {
//			newHigh = append(newHigh, dataHigh[i])
//			newLow = append(newLow, dataLow[i])
//		}
//	}
//	return newHigh, newLow
//}
//
//func isPen(data []float64) bool {
//	if len(data) < 3 {
//		return false
//	}
//	return true
//}
//func isTopFractal(data []float64) bool {
//	return false
//}
//func isBottomFractal(data []float64) bool {
//	return false
//}
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

		ac := &avgContext{
			State:          STATE_UNKOWN,
			profit:         0.0,
			Sell_Min_Value: -1,
			Buy_Max_Value:  -1,
		}

		tmpArr := make([]float64, 0)
		//for i := 0; i < 500; i++ {
		//	getWave(&stockBig, i)
		//}
		for i := 1; i < len(stockBig.dataClose); i += 1 {
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

			//drawAllMinMax(p, &stockBig, 2, black)
			//drawMinMax(p, stockBig.dataHigh[start:end], stockBig.dataMinMax[start:end], 1, 3, black)
			//drawMinMax(p, stockBig.dataLow[start:end], stockBig.dataMinMax[start:end], -1, 3, dark_red)

			st := copyStock(&stockBig, start, i+1)
			arr, _ := getAllRect(st)
			if len(arr) > 0 {
				for _, r := range arr {
					drawRectangle(p, r.left, r.top, r.right, r.bottom, gray)
				}
			}
			drawAllMinMax(p, st, 2, black)

			filename := fmt.Sprintf("/Users/xinmei365/stock/%03d_%03d.png", index, i)
			ret := run(ac, p, st, filename, i)
			if ret {
				tmpArr = append(tmpArr, ac.profit)
				log.Infof("[%d] profit:%f", i, ac.profit)
			}
			//p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
		}
		//
		drawPic(tmpArr, "Count", "Profit", fmt.Sprintf("/Users/xinmei365/profilt_%d.png", index))
	}
}
