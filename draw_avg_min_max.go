package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
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
	// Make a line plotter and set its style.
	//l, err := plotter.NewLine(points)
	//if err != nil {
	//	panic(err)
	//}
	//l.LineStyle.Width = vg.Points(1)
	//l.LineStyle.Width = vg.Points(1)
	//l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
	//l.LineStyle.Color = clr

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
		if stock.dataMinMax[i] == 2 || stock.dataMinMax[i] == 2 {
			var elem = plotter.XY{
				float64(i), stock.dataHigh[i],
			}
			points = append(points, elem)
		}
		if stock.dataMinMax[i] == -2 || stock.dataMinMax[i] == -2 {
			var elem = plotter.XY{
				float64(i), stock.dataLow[i],
			}
			points = append(points, elem)
		}
	}
	// Make a line plotter and set its style.
	//l, err := plotter.NewLine(points)
	//if err != nil {
	//	panic(err)
	//}
	//l.LineStyle.Width = vg.Points(1)
	//l.LineStyle.Width = vg.Points(1)
	//l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
	//l.LineStyle.Color = clr

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
	// Make a line plotter and set its style.
	//l, err := plotter.NewLine(points)
	//if err != nil {
	//	panic(err)
	//}
	//l.LineStyle.Width = vg.Points(1)
	//l.LineStyle.Width = vg.Points(1)
	//l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
	//l.LineStyle.Color = clr

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

	//i := rand.Intn(len(colorArray))
	//if flag == 0 {
	//	lpLine.LineStyle.Color = black
	//} else {
	//	lpLine.LineStyle.Color = red
	//}
	lpLine.LineStyle.Color = clr
	//flag = (flag + 1) % 2
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

//func drawPoint3(p *plot.Plot, x1 float64, y1 float64, radius int, clr color.Color) {
//	points := plotter.XYs{
//		{x1, y1},
//	}
//	plotutil.AddLinePoints(p, points)
//	_, lpPoint, err := plotter.NewLinePoints(points)
//	if err != nil {
//		panic(err)
//	}
//	lpPoint.Radius = vg.Length(radius)
//	lpPoint.Color = clr
//	p.Add(lpPoint)
//}
//func drawPoint(p *plot.Plot, x1 float64, y1 float64, radius int) {
//	points := plotter.XYs{
//		{x1, y1},
//	}
//	plotutil.AddLinePoints(p, points)
//	_, lpPoint, err := plotter.NewLinePoints(points)
//	if err != nil {
//		panic(err)
//	}
//	lpPoint.Radius = vg.Length(radius)
//	lpPoint.Color = blue
//	p.Add(lpPoint)
//}
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

//func drawAvg30Min(p *plot.Plot) {
//	points := make(plotter.XYs, 0)
//	for i, val := range stock.avg30MinMax {
//		var x float64
//		x = float64(i)
//		if val == -1 {
//			var elem = plotter.XY{
//				x, stock.avg30[val],
//			}
//			points = append(points, elem)
//		}
//	}
//	// Make a line plotter and set its style.
//	l, err := plotter.NewLine(points)
//	if err != nil {
//		panic(err)
//	}
//	l.LineStyle.Width = vg.Points(1)
//	l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
//	l.LineStyle.Color = color.RGBA{B: 255, A: 255}
//
//	p.Add(l)
//}
//
//func drawAvg30Max(p *plot.Plot) {
//	points := make(plotter.XYs, 0)
//	for i, flag := range stock.avg30MinMax {
//		var x float64
//		x = float64(i)
//		if flag == 1 {
//			var elem = plotter.XY{
//				x, stock.avg30[i],
//			}
//			points = append(points, elem)
//		}
//	}
//	// Make a line plotter and set its style.
//	l, err := plotter.NewLine(points)
//	if err != nil {
//		panic(err)
//	}
//	l.LineStyle.Width = vg.Points(1)
//	l.LineStyle.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}
//	l.LineStyle.Color = color.RGBA{B: 255, A: 255}
//
//	lpLine, lpPoints, err := plotter.NewLinePoints(points)
//	if err != nil {
//		panic(err)
//	}
//	lpLine.Color = color.RGBA{G: 255, A: 255}
//	lpPoints.Shape = draw.CircleGlyph{}
//	lpPoints.Radius = 2
//	lpPoints.Color = color.RGBA{R: 255, A: 255}
//
//	p.Add(l, lpPoints)
//
//	//p.Add(l)
//}
//
//func drawAvg150MinMax(p *plot.Plot) {
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
