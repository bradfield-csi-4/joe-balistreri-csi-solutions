package node

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"
)

type ScanNode struct {
	reader      *bufio.Reader
	initialized bool
	columnNames []string
}

var tablenameToFilename = map[string]string{
	"movies":        "data/ml-20m/movies.csv",
	"ratings":       "data/ml-20m/ratings.csv",
	"tags":          "data/ml-20m/tags.csv",
	"links":         "data/ml-20m/links.csv",
	"genome-scores": "data/ml-20m/genome-scores.csv",
	"genome-tags":   "data/ml-20m/genome-tags.csv",
}

func NewScanNode(tablename string) *ScanNode {
	filename, ok := tablenameToFilename[tablename]
	if !ok {
		panic("invalid tablename")
	}
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return &ScanNode{reader: bufio.NewReader(f)}
}

func (s *ScanNode) Next() map[string]string {
	if !s.initialized {
		header := s.read()
		if header == nil {
			return nil
		}
		s.columnNames = s.readCsv(*header)
		s.initialized = true
	}

	nextData := s.read()
	if nextData == nil {
		return nil
	}
	result := map[string]string{}
	csvData := s.readCsv(*nextData)
	for i := range csvData {
		result[s.columnNames[i]] = csvData[i]
	}
	return result
}

func (s *ScanNode) readCsv(input string) []string {
	res, err := csv.NewReader(strings.NewReader(input)).Read()
	if err != nil {
		panic(err)
	}
	return res
}

func (s *ScanNode) read() *string {
	line, _, err := s.reader.ReadLine()
	if err != nil {
		return nil
	}
	lString := string(line)
	return &lString
}
