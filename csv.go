package main

import (
	"encoding/csv"
	"io"
)

func readData(r io.Reader, offset, limit int) ([]map[string]interface{}, error) {
	ms := make([]map[string]interface{}, 0)
	csvr := csv.NewReader(r)
	headers, err := csvr.Read()
	if err != nil {
		return nil, err
	}

	var rowIndex int
	for {
		rowIndex++
		if limit > 0 && rowIndex >= 1+offset+limit {
			break
		}
		if rowIndex <= offset {
			continue
		}

		row, err := csvr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i := range row {
			m[headers[i]] = row[i]
		}
		ms = append(ms, m)
	}

	return ms, nil
}
