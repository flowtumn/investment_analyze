package main

import (
	"bufio"
	"fmt"
	"os"
)

type AnalyzeResult struct {
	TotalCount   int64
	HistoryCount map[int]int64
	HistoryData  map[int]float64
}

func (self *AnalyzeResult) Initialize() {
	self.TotalCount = 0
	self.HistoryCount = map[int]int64{}
	self.HistoryData = map[int]float64{}
}

func (self *AnalyzeResult) AddData(yesterdayDate, newDate InvestmentElemDate) {
	self.TotalCount = self.TotalCount + 1
	self.HistoryCount[newDate.Date.Day] = self.HistoryCount[newDate.Date.Day] + 1
	self.HistoryData[newDate.Date.Day] = self.HistoryData[newDate.Date.Day] + newDate.Price.Close - yesterdayDate.Price.Close
}

func (self *AnalyzeResult) Dump() {
	fmt.Printf("TotalSamples: %d\n", self.TotalCount)
	for i := 1; i <= 31; i++ {
		fmt.Printf("Day %02d: Samples: %d   Score: %f\n", i, self.HistoryCount[i], self.HistoryData[i])
	}
}

func analyze(csvPath string, reverse bool, parser CsvParser) (*AnalyzeResult, error) {
	fp, err := os.Open(csvPath)
	if nil != err {
		return nil, err
	}

	defer func() {
		fp.Close()
	}()

	r := AnalyzeResult{}
	r.Initialize()

	scanner := bufio.NewScanner(fp)

	//Get a first data.
	for scanner.Scan() {
		baseData, err := parser.ParseLine(scanner.Text())
		if nil == err && nil != baseData {
			//analyze.
			for scanner.Scan() {
				data, err := parser.ParseLine(scanner.Text())
				if nil == err && nil != data {
					if reverse {
						r.AddData(*data, *baseData)
					} else {
						r.AddData(*baseData, *data)
					}
					baseData = data
				}
			}
			break
		}
	}

	return &r, nil
}

func main() {
	result, err := analyze(os.Args[1], true, &NikkeiParser{})
	if nil != err {
		panic(err)
	}

	result.Dump()
}
