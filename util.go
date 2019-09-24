package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"math"
)

func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func getVariance(data []float64, start int, end int) float64 {
	if start >= end {
		log.Errorf("error!")
		return -1
	}
	avg := getAvg(data, start, end)
	variance := 0.0
	for i := start; i < end; i++ {
		tmp := (data[i] - avg)
		tmp02 := tmp * tmp
		variance += tmp02
	}
	cnt := float64(end - start)
	variance /= cnt
	return variance
}

func getSum(data []float64, start int, end int) float64 {
	if start >= end {
		log.Errorf("error!")
		return -1
	}
	total := 0.0
	for i := start; i < end; i++ {
		total += data[i]
	}
	return total
}

func getAvg(data []float64, start int, end int) float64 {
	if start >= end {
		log.Errorf("error!")
		return -1
	}
	total := 0.0
	for i := start; i < end; i++ {
		total += data[i]
	}
	cnt := end - start
	return total / float64(cnt)
}

func minimum(data []float64) float64 {
	if data == nil {
		return math.NaN()
	}
	min := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] < min {
			min = data[i]
		}
	}
	return min
}

func maximum(data []float64) float64 {
	if data == nil {
		return math.NaN()
	}
	max := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] > max {
			max = data[i]
		}
	}
	return max
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func isign(a int) int {
	if a < 0 {
		return -1
	}
	if a == 0 {
		return 0
	}
	return 1
}

func clip(x, l, u int) int {
	if x < imin(l, u) {
		return l
	}
	if x > imax(l, u) {
		return u
	}
	return x
}

func fmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
