package main

import (
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

func run(ac *avgContext, p *plot.Plot, data []float64, filename string, pos int) bool {

	_, st := getAllRect(data)
	curPos := len(data) - 1

	//log.Infof("len(arr):%v", len(arr))
	//for _, r := range arr {
	//	//drawRectangle(p, r.left, r.top, r.right, r.bottom, gray)
	//	//drawLine(p, r.left, r.top, r.right, r.bottom)
	//	//drawLine(p, r.left, r.bottom, r.right, r.top)
	//}
	//drawMinMax(p, st.dataClose[0:pos], st.dataMinMax[0:pos], 1, 1, yellow)
	//p.Save(vg.Length(picwidth), vg.Length(picheight), filename)

	if ac.State == STATE_UNKOWN {
		ok, revert, change := isValidInit(ac, st.dataClose)
		if ok && revert && change {
			//log.Infof("ac: %s", ac.Show())

			drawPoint2(p, float64(curPos), st.dataClose[curPos], 20, red)
			if ac.Action == ACTION_BUY {
				drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
			} else if ac.Action == ACTION_SELL {
				drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
			}

			p.X.Label.Text = ac.State + " " + ac.Action
			p.Save(vg.Length(picwidth), vg.Length(picheight), filename)

			return true
		}
	} else if ac.State != STATE_UNKOWN {
		ok, revert, change, arr := forwardState(ac, st.dataClose)
		for _, r := range arr {
			drawRectangle(p, r.left, r.top, r.right, r.bottom, yellow)
		}
		drawData(p, st.dense, 1, gray)
		//log.Infof("pos:%d ok:%v", pos, ok)

		if ok {
			p.X.Label.Text = ac.State + " " + ac.Action

			red := exampleThumbnailer{Color: color.NRGBA{R: 255, A: 255}}
			green := exampleThumbnailer{Color: color.NRGBA{G: 255, A: 255}}
			blue := exampleThumbnailer{Color: color.NRGBA{B: 255, A: 255}}

			l, err := plot.NewLegend()
			if err != nil {
				panic(err)
			}
			l.Add("red", red)

			if revert {
				//log.Infof("ac: %s", ac.Show())

				drawPoint2(p, float64(curPos), st.dataClose[curPos], 20, black)

				if ac.Action == ACTION_BUY {
					drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
				} else if ac.Action == ACTION_SELL {
					drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
				}
				p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
				return true
			} else if change {
				//log.Infof("ac: %s", ac.Show())

				drawPoint2(p, float64(curPos), st.dataClose[curPos], 20, black)

				if ac.Action == ACTION_BUY {
					drawRectangle(p, ac.Buy_stop.left, ac.Buy_stop.top, ac.Buy_stop.right, ac.Buy_stop.bottom, blue)
				} else if ac.Action == ACTION_SELL {
					drawRectangle(p, ac.Sell_stop.left, ac.Sell_stop.top, ac.Sell_stop.right, ac.Sell_stop.bottom, green)
				}
				p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
				return true
			}
		}
	} else {
		log.Infof("ignore:%d", pos)
	}
	return false
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
		//stock.GetDist()

		ac := &avgContext{
			State:  STATE_UNKOWN,
			profit: 0.0,
		}

		tmpArr := make([]float64, 0)
		for i := 1; i < len(stock.dataClose); i += 1 {
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

			//1， 绘制底图
			drawData(p, stock.dataClose[start:end], 1, pink)
			//drawData(p, stock.dense[start:end], 1, orange)
			//drawMinMax(p, stock.dataClose[start:end], stock.dataMinMax[start:end], 1, 3, gray)
			//drawMinMax(p, stock.dataClose[start:end], stock.dataMinMax[start:end], -1, 3, gray)
			filename := fmt.Sprintf("/Users/xinmei365/stock/%03d_%03d.png", index, i)
			//arr, _ := getAllRect(stock.dataClose[start : i+1])
			//if len(arr) > 0 {
			//	for _, r := range arr {
			//		drawRectangle(p, r.left, r.top, r.right, r.bottom, gray)
			//	}
			//}
			ret := run(ac, p, stock.dataClose[start:i+1], filename, i)
			if ret {
				log.Infof("[%d] profit:%f", i, ac.profit)
				tmpArr = append(tmpArr, ac.profit)
			}
			//p.Save(vg.Length(picwidth), vg.Length(picheight), filename)
		}
		//
		drawPic(tmpArr, "Count", "Profit", "/Users/xinmei365/profilt.png")
	}
}
