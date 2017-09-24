package main

import "strconv"

func ToFloat64(str string, def float64) float64 {
	v, err := strconv.ParseFloat(str, 64)
	if nil != err {
		return def
	}
	return v
}

func ToInt(str string, def int) int {
	v, err := strconv.Atoi(str)
	if nil != err {
		return def
	}
	return v
}
