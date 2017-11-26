package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type AnalyzeData struct {
	Count               int64
	TotalCost			int64
	Score               float64
	PreviousHigherCount int64
	PreviousSameCount   int64
	PreviousLowerCount  int64
}

func (self AnalyzeData) ToCsv() string {
	//sample数, TotalCost, Costの平均値, Score(その日に預けていたらの増減), 前日より高値の回数,前日より低値の回数,前日と同値の回数
	return fmt.Sprintf("%d,%d,%d,%g,%d,%d,%d", self.Count, self.TotalCost, self.TotalCost / self.Count, self.Score, self.PreviousHigherCount, self.PreviousLowerCount, self.PreviousSameCount)
}

type AnalyzeResult struct {
	TotalCount   int64
	HistoryDatas map[int]interface{}
}

func (self *AnalyzeResult) Initialize() {
	self.TotalCount = 0
	self.HistoryDatas = map[int]interface{}{}

	for i := 1; i <= 31; i++ {
		self.HistoryDatas[i] = AnalyzeData{}
	}
}

func (self *AnalyzeResult) AddData(yesterdayData, newData InvestmentElemData) {
	self.TotalCount = self.TotalCount + 1
	v, ok := self.HistoryDatas[newData.Date.Day].(AnalyzeData)
	if !ok {
		return
	}

	v.Count = v.Count + 1
	v.TotalCost = v.TotalCost + (int64)(newData.Price.Close)

	//本日が前日に比べて高値で終わったのかを計算。
	v.Score = newData.Price.Close - yesterdayData.Price.Close

	//前日より高値、低値かのカウントアップ
	if yesterdayData.Price.Close < newData.Price.Close {
		v.PreviousHigherCount = v.PreviousHigherCount + 1
	} else if yesterdayData.Price.Close > newData.Price.Close {
		v.PreviousLowerCount = v.PreviousLowerCount + 1
	} else {
		v.PreviousSameCount = v.PreviousSameCount + 1
	}

	self.HistoryDatas[newData.Date.Day] = v
}

func (self *AnalyzeResult) DumpCSV(path string) error {
	fp, err := os.Create(path)
	if nil != err {
		return err
	}

	defer func() {
		fp.Close()
	}()

	for i := 1; i <= 31; i++ {
		v, ok := self.HistoryDatas[i].(AnalyzeData)
		if !ok {
			continue
		}
		fp.Write([]byte(fmt.Sprintf("%02d,%s\n", i, v.ToCsv())))
	}

	fp.Write([]byte(fmt.Sprintf("TotalSamples: %d", self.TotalCount)))
	return nil
}

func analyze(csvPath string, parser CsvParser) (*AnalyzeResult, error) {
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

	samples := []InvestmentElemData{}

	//Parse Data.
	for scanner.Scan() {
		data, err := parser.ParseLine(scanner.Text())
		if nil == err && nil != data {
			samples = append(samples, *data)
		}
	}

	//Sort.
	sort.Slice(samples, func(i, j int) bool {
		lhs := fmt.Sprintf("%04d/%02d/%02d", samples[i].Date.Year, samples[i].Date.Month, samples[i].Date.Day)
		rhs := fmt.Sprintf("%04d/%02d/%02d", samples[j].Date.Year, samples[j].Date.Month, samples[j].Date.Day)
		return strings.Compare(lhs, rhs) < 0
	})

	baseData := samples[0]
	for _, v := range samples[1:] {
		r.AddData(baseData, v)
		baseData = v
	}

	return &r, nil
}

func main() {
	for i := 1 ; i < len(os.Args); i++ {
		result, err := analyze(os.Args[i], &TrustParser{})
		if nil == err {
			savePath := os.Args[i] + "_analyze.csv"
			result.DumpCSV(savePath)
		} else {
			fmt.Printf("analyze error: %s\n", err.Error())
		}
	}
}
