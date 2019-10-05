package main

// 根据前面的 length 个值，获取平均值
func get_pre_avg(value_list []float64, index int, length int) float64 {
	var total float64
	var count int
	total = 0.0
	count = 0
	for i := 0; i < length; i++ {
		pos := index - i
		if pos < 0 {
			break
		}
		total += value_list[pos]
		count += 1
	}
	return total / float64(count)
}
