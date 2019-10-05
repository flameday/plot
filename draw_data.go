package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

func drawData2(p *plot.Plot, data []int, width float64, clr color.Color) {
	points := make(plotter.XYs, 0)
	for i, val := range data {
		var x float64
		x = float64(i)
		var elem = plotter.XY{
			x, float64(val),
		}
		points = append(points, elem)
	}

	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpLine.Color = clr
	lpPoints.Shape = draw.CircleGlyph{}
	lpPoints.Radius = vg.Length(width)
	lpPoints.Color = clr

	p.Add(lpLine, lpPoints)
}

func drawData(p *plot.Plot, data []float64, width float64, clr color.Color) {
	points := make(plotter.XYs, 0)
	for i, val := range data {
		var x float64
		x = float64(i)
		var elem = plotter.XY{
			x, val,
		}
		points = append(points, elem)
	}

	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpLine.Color = clr
	lpPoints.Shape = draw.CircleGlyph{}
	lpPoints.Radius = vg.Length(width)
	lpPoints.Color = clr

	p.Add(lpLine, lpPoints)
}
