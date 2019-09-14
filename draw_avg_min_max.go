package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

func drawAvg30Min(p *plot.Plot) {
	points := make(plotter.XYs, 0)
	for i, val := range stock.avg30MinMax {
		var x float64
		x = float64(i)
		if val == -1 {
			var elem = plotter.XY{
				x, stock.avg30[val],
			}
			points = append(points, elem)
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

	p.Add(l)
}

func drawAvg30Max(p *plot.Plot) {
	points := make(plotter.XYs, 0)
	for i, flag := range stock.avg30MinMax {
		var x float64
		x = float64(i)
		if flag == 1 {
			var elem = plotter.XY{
				x, stock.avg30[i],
			}
			points = append(points, elem)
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

	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpLine.Color = color.RGBA{G: 255, A: 255}
	lpPoints.Shape = draw.CircleGlyph{}
	lpPoints.Radius = 2
	lpPoints.Color = color.RGBA{R: 255, A: 255}

	p.Add(l, lpPoints)

	//p.Add(l)
}

func drawAvg150MinMax(p *plot.Plot) {
	points := make(plotter.XYs, 0)
	for i, val := range stock.avg150 {
		var x float64
		x = float64(i)
		var elem = plotter.XY{
			x, val,
		}
		points = append(points, elem)
	}
	// Make a line plotter and set its style.
	l, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
	l.LineStyle.Color = color.RGBA{R: 255, A: 255}

	p.Add(l)
}
