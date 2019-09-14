package main

import (
	"github.com/prometheus/common/log"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

var (
	white     color.Color = color.RGBA{255, 255, 255, 255}
	blue      color.Color = color.RGBA{0, 0, 255, 255}
	red       color.Color = color.RGBA{255, 0, 0, 255}
	picwidth  float64     = 512 * 1
	picheight float64     = 384 * 1
	stock     Stock
)

// 原始的点
func drawCircle(p *plot.Plot) {
	points := make(plotter.XYs, 0)
	for i, val := range stock.data {
		var x float64
		x = float64(i)
		var elem = plotter.XY{
			x, val,
		}
		points = append(points, elem)
		if i < 10 {
			log.Info("[%d] %v", i, elem)
		}
	}
	// Make a line plotter and set its style.
	l, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	//圈圈
	s, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	//线段
	// Make a line plotter with points and set its style.
	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpLine.Color = color.RGBA{G: 255, A: 255}
	lpPoints.Shape = draw.CircleGlyph{}
	lpPoints.Radius = 1
	lpPoints.Color = color.RGBA{R: 255, A: 255}

	p.Add(lpPoints)
}
