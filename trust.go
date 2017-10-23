package main

import (
	"encoding/csv"
	"errors"
	"regexp"
	"strings"
)

type TrustParser struct {
	CsvParser
}

func (self *TrustParser) ConvertDate(date string) *Date {
	r := regexp.MustCompile(`([0-9]{4})\/([0-9]{2})\/([0-9]{2})`)
	data := r.FindSubmatch([]byte(date))

	if 4 != len(data) {
		return nil
	}

	return &Date{
		Year:  ToInt(string(data[1]), 0),
		Month: ToInt(string(data[2]), 0),
		Day:   ToInt(string(data[3]), 0),
	}
}

func (self *TrustParser) ParseLine(csvLine string) (*InvestmentElemData, error) {
	parser := csv.NewReader(strings.NewReader(csvLine))
	record, err := parser.Read()
	if nil != err {
		return nil, err
	}

	if 4 > len(record) {
		return nil, errors.New("parser.Read error.")
	}

	date := self.ConvertDate(record[0])
	if nil == date {
		return nil, errors.New("ConvertDate error.")
	}

	return &InvestmentElemData{
		Date: *date,
		Price: Price{
			Open:  ToFloat64(record[1], 0.0),
			High:  ToFloat64(record[1], 0.0),
			Low:   ToFloat64(record[1], 0.0),
			Close: ToFloat64(record[1], 0.0),
		},
	}, nil
}
