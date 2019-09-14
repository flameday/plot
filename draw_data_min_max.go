package main

//func drawMax(p *plot.Plot) {
//	points := make(plotter.XYs, 0)
//	for i, flag := range stock.dataMinMax {
//		if flag == 1 {
//			var elem = plotter.XY{
//				float64(i), stock.data[i],
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
//	l.LineStyle.Color = color.RGBA{R: 128, G: 255, A: 255}
//
//	//圈圈
//	s, err := plotter.NewScatter(points)
//	if err != nil {
//		panic(err)
//	}
//	s.GlyphStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
//	p.Add(l, s)
//}
//
//func drawMin(p *plot.Plot) {
//	points := make(plotter.XYs, 0)
//	for i, flag := range stock.dataMinMax {
//		if flag == -1 {
//			var elem = plotter.XY{
//				float64(i), stock.data[i],
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
//	l.LineStyle.Color = color.RGBA{R: 2, G: 255, B: 255, A: 255}
//
//	//圈圈
//	s, err := plotter.NewScatter(points)
//	if err != nil {
//		panic(err)
//	}
//	s.GlyphStyle.Color = color.RGBA{R: 2, G: 255, B: 255, A: 255}
//	p.Add(l, s)
//}
