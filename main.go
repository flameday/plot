package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"image/color"
	"time"
)

var (
	white     color.Color = color.RGBA{255, 255, 255, 255}
	blue      color.Color = color.RGBA{0, 0, 255, 255}
	red       color.Color = color.RGBA{255, 0, 0, 255}
	picwidth  float64     = 512 * 1
	picheight float64     = 384 * 1
	stock     Stock
)

// 大家可以查看这个网址看看这个image包的使用方法 http://golang.org/doc/articles/image_draw.html
func main() {
	stock.LoadData()

	//创建 plog
	p, _ := plot.New()
	t := time.Now()

	p.Title.Text = t.Format("2006-01-02 15:04:05.000000000")
	p.X.Label.Text = "Quantity Demand"
	p.Y.Label.Text = "Price"

	//画圈圈
	//drawCircle(p)
	//画平均线
	//drawAvg150(p)
	//drawAvg30(p)
	drawData(p, stock.data, 1, blue)
	drawData(p, stock.avgMiddle, 1, color.RGBA{R: 255, G: 128, B: 255, A: 255})
	drawMinMax(p, stock.avgMiddle, stock.avgMiddleMinMax, 1, 3, color.RGBA{R: 128, G: 128, B: 255, A: 255})
	drawMinMax(p, stock.avgMiddle, stock.avgMiddleMinMax, -1, 1, color.RGBA{R: 128, G: 128, B: 255, A: 255})
	drawData(p, stock.avg30, 1, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	//drawAvg30Max(p)
	//drawMax(p)
	//drawMin(p)
	//drawAvg30Max(p)
	//drawCleanMin(p)
	//p.Add(s, l, lpLine, lpPoints)
	//plotutil.AddLinePoints(p, points)

	p.Save(vg.Length(picwidth), vg.Length(picheight), "/Users/xinmei365/price.png")

	//stock.LoadData("/Users/xinmei365/stock_data_history/day/data/000002.csv")
	//http.HandleFunc("/", RrawPicture)
	//http.ListenAndServe(":999", nil)
}
