package main

import (
	"encoding/csv"
	"errors"
	"regexp"
	"strings"
)

type NikkeiParser struct {
	CsvParser
}

func (self *NikkeiParser) ConvertDate(date string) *Date {
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

func (self *NikkeiParser) ParseLine(csvLine string) (*InvestmentElemDate, error) {
	parser := csv.NewReader(strings.NewReader(csvLine))
	record, err := parser.Read()
	if nil != err {
		return nil, err
	}

	date := self.ConvertDate(record[0])
	if nil == date {
		return nil, errors.New("ConvertDate error.")
	}

	return &InvestmentElemDate{
		Date: *date,
		Price: Price{
			Open:  ToFloat64(record[1], 0.0),
			High:  ToFloat64(record[2], 0.0),
			Low:   ToFloat64(record[3], 0.0),
			Close: ToFloat64(record[4], 0.0),
		},
	}, nil
}
