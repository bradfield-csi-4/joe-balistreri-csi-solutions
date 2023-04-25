package node

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"os"
	"strings"
)

type ScanNode struct {
	reader      *bufio.Reader
	initialized bool
}

var tablenameToFilename = map[string]string{
	"movies":        "data/ml-20m/movies.csv",
	"ratings":       "data/ml-20m/ratings.csv",
	"tags":          "data/ml-20m/tags.csv",
	"links":         "data/ml-20m/links.csv",
	"genome-scores": "data/ml-20m/genome-scores.csv",
	"genome-tags":   "data/ml-20m/genome-tags.csv",
}

func NewTestScanNode(lines []string) *ScanNode {
	return &ScanNode{reader: bufio.NewReader(bytes.NewReader([]byte(strings.Join(lines, "\n"))))}
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

func (s *ScanNode) Next() Row {
	nextData := s.read()
	if nextData == nil {
		return nil
	}
	return s.readCsv(*nextData)
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
