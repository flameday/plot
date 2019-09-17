package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

//func drawData2(p *plot.Plot, dataClose []int, width float64, clr color.Color) {
//	points := make(plotter.XYs, 0)
//	for i, val := range dataClose {
//		var x float64
//		x = float64(i)
//		var elem = plotter.XY{
//			x, float64(val),
//		}
//		points = append(points, elem)
//	}
//	// Make a line plotter and set its style.
//	//l, err := plotter.NewLine(points)
//	//if err != nil {
//	//	panic(err)
//	//}
//	////l.LineStyle.Width = vg.Points(1)
//	//l.LineStyle.Width = vg.Points(width)
//	//l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
//	//l.LineStyle.Color = clr
//	//l.LineStyle.Color = color.RGBA{R: 255, B: 255, A: 255}
//
//	//p.Add(l)
//
//	lpLine, lpPoints, err := plotter.NewLinePoints(points)
//	if err != nil {
//		panic(err)
//	}
//	lpLine.Color = clr
//	lpPoints.Shape = draw.CircleGlyph{}
//	lpPoints.Radius = vg.Length(width)
//	lpPoints.Color = clr
//
//	p.Add(lpLine, lpPoints)
//}

func drawData(p *plot.Plot, dataClose []float64, width float64, clr color.Color) {
	points := make(plotter.XYs, 0)
	for i, val := range dataClose {
		var x float64
		x = float64(i)
		var elem = plotter.XY{
			x, val,
		}
		points = append(points, elem)
	}
	// Make a line plotter and set its style.
	//l, err := plotter.NewLine(points)
	//if err != nil {
	//	panic(err)
	//}
	////l.LineStyle.Width = vg.Points(1)
	//l.LineStyle.Width = vg.Points(width)
	//l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
	//l.LineStyle.Color = clr
	//l.LineStyle.Color = color.RGBA{R: 255, B: 255, A: 255}

	//p.Add(l)

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

//func drawAvg30(p *plot.Plot) {
//	points := make(plotter.XYs, 0)
//	for i, val := range stock.avg30 {
//		var x float64
//		x = float64(i)
//		var elem = plotter.XY{
//			x, val,
//		}
//		points = append(points, elem)
//	}
//	// Make a line plotter and set its style.
//	l, err := plotter.NewLine(points)
//	if err != nil {
//		panic(err)
//	}
//	l.LineStyle.Width = vg.Points(1)
//	l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
//	l.LineStyle.Color = color.RGBA{R: 255, B: 255, A: 255}
//
//	p.Add(l)
//}
//func drawAvg150(p *plot.Plot) {
//	points := make(plotter.XYs, 0)
//	for i, val := range stock.avg150 {
//		var x float64
//		x = float64(i)
//		var elem = plotter.XY{
//			x, val,
//		}
//		points = append(points, elem)
//	}
//	// Make a line plotter and set its style.
//	l, err := plotter.NewLine(points)
//	if err != nil {
//		panic(err)
//	}
//	l.LineStyle.Width = vg.Points(1)
//	l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
//	l.LineStyle.Color = color.RGBA{R: 255, A: 255}
//
//	p.Add(l)
//}
