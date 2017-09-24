package main

import "fmt"

func analyze(parser CsvParser) {
	v, _ := parser.ParseLine("2017/08/10,19792.45,19829.88,19685.83,19729.74")

	record := map[int]float64{}

	fmt.Printf("%+v\n", *v)
	//CSVをParse。
	//前日の終値と、今回の終値の差を記録。
	currentClose := 100.2
	record[v.Date.Day] = record[v.Date.Day] + (v.Price.Close - currentClose)
}

func main() {
	analyze(&NikkeiParser{})
	analyze(&TopixParser{})
}
