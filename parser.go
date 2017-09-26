package main

type Date struct {
	Year  int
	Month int
	Day   int
}

type Price struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

type InvestmentElemData struct {
	Date  Date
	Price Price
}

type CsvParser interface {
	ParseLine(csvLine string) (*InvestmentElemData, error)
}
