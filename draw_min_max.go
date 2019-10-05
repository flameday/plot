package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

func drawMinMax(p *plot.Plot, data []float64, posArr []int, minMax int, width float64, clr color.Color) {
	points := make(plotter.XYs, 0)
	for i, val := range data {
		var x float64
		x = float64(i)
		if posArr[i] == minMax {
			var elem = plotter.XY{
				x, val,
			}
			points = append(points, elem)
		}
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

func drawAllSubMinMax(p *plot.Plot, stock *Stock, width float64, clr color.Color) {
	points := make(plotter.XYs, 0)
	for i, _ := range stock.dataClose {
		if stock.subDataMinMax[i] > 0 || stock.dataMinMax[i] > 0 {
			var elem = plotter.XY{
				float64(i), stock.dataHigh[i],
			}
			points = append(points, elem)
		}
		if stock.subDataMinMax[i] < 0 || stock.dataMinMax[i] < 0 {
			var elem = plotter.XY{
				float64(i), stock.dataLow[i],
			}
			points = append(points, elem)
		}
	}

	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpLine.Color = clr
	lpPoints.Shape = draw.CircleGlyph{}
	lpPoints.Radius = vg.Length(width)
	lpPoints.Color = clr
	lpLine.LineStyle.Width = vg.Points(2)

	p.Add(lpLine, lpPoints)
}

func drawAllMinMax(p *plot.Plot, stock *Stock, width float64, clr color.Color) {
	points := make(plotter.XYs, 0)
	for i, _ := range stock.dataClose {
		if stock.dataMinMax[i] == 1 {
			var elem = plotter.XY{
				float64(i), stock.dataHigh[i],
			}
			points = append(points, elem)
		}
		if stock.dataMinMax[i] == -1 {
			var elem = plotter.XY{
				float64(i), stock.dataLow[i],
			}
			points = append(points, elem)
		}
	}

	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}
	lpLine.Color = clr
	lpPoints.Shape = draw.CircleGlyph{}
	lpPoints.Radius = vg.Length(width)
	lpPoints.Color = clr
	lpLine.LineStyle.Width = vg.Points(2)

	p.Add(lpLine, lpPoints)
}
