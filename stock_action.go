package main

func action() {

	for i := 0; i < 10000; i++ {
		posMax1 := findNextIndex(stock.resetMinMax, 0, MAX_VALUE_FLAG)
		posMin1 := findNextIndex(stock.resetMinMax, posMax1+1, MIN_VALUE_FLAG)
		posMax2 := findNextIndex(stock.resetMinMax, posMin1+1, MAX_VALUE_FLAG)
		posMin2 := findNextIndex(stock.resetMinMax, posMax2+1, MIN_VALUE_FLAG)
		if posMax1 == -1 || posMin1 == -1 || posMax2 == -1 || posMin2 == -1 {
			break
		}

	}

}
