package main

//func (stock *Stock) GetMacd() {
//	for i := 1; i < len(stock.dataClose); i++ {
//		ema12 := (stock.dataClose[i-1]*11 + stock.dataClose[i]*2) / 13.0
//		ema26 := (stock.dataClose[i-1]*25 + stock.dataClose[i]*2) / 27.0
//		stock.DIFF[i] = ema12 - ema26
//		stock.DEA[i] = stock.DEA[i-1]*8/10.0 + stock.DIFF[i]*2/10.0
//		stock.BAR[i] = 2 * (stock.DIFF[i] - stock.DEA[i])
//	}
//}
