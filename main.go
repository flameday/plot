package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"time"
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
	drawCircle(p)
	//画平均线
	//drawAvg150(p)
	drawAvg30(p)
	//drawAvg30Max(p)
	drawMax(p)
	drawMin(p)
	drawAvg30Max(p)
	//drawCleanMin(p)
	//p.Add(s, l, lpLine, lpPoints)
	//plotutil.AddLinePoints(p, points)

	p.Save(vg.Length(picwidth), vg.Length(picheight), "/Users/xinmei365/price.png")

	//stock.LoadData("/Users/xinmei365/stock_data_history/day/data/000002.csv")
	//http.HandleFunc("/", RrawPicture)
	//http.ListenAndServe(":999", nil)
}
