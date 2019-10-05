package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"image/color"
)

func drawRectangle(p *plot.Plot, x1 float64, y1 float64, x2 float64, y2 float64, clr color.Color) {
	points := plotter.XYs{
		{x1, y1},
		{x1, y2},
		{x2, y2},
		{x2, y1},
		{x1, y1},
	}
	plotutil.AddLinePoints(p, points)
	lpLine, _, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpLine.LineStyle.Width = vg.Points(2)

	lpLine.LineStyle.Color = clr
	p.Add(lpLine)
}

func drawPoint(p *plot.Plot, x1 float64, y1 float64, radius int, clr color.Color) {
	points := plotter.XYs{
		{x1, y1},
	}
	//plotutil.AddLinePoints(p, points)
	_, lpPoint, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpPoint.Radius = vg.Length(radius)
	lpPoint.Color = clr
	p.Add(lpPoint)
}

func drawLine(p *plot.Plot, x1 float64, y1 float64, x2 float64, y2 float64) {
	points := plotter.XYs{
		{x1, y1},
		{x2, y2},
	}
	plotutil.AddLinePoints(p, points)
	lpLine2, _, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpLine2.Color = red
	p.Add(lpLine2)
}
