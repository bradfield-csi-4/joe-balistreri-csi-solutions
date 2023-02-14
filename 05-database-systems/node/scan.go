package node

import (
	"encoding/csv"
	"strings"
)

type ScanNode struct {
	index       int
	input       []string
	initialized bool
	columnNames []string
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
	if s.index >= len(s.input) {
		return nil
	}
	res := &s.input[s.index]
	s.index++
	return res
}
