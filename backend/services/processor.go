package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// ProcessCSV reads a CSV from the provided reader and sums sales by city.
func ProcessCSV(file io.Reader) (*map[string]float64, error) {
	reader := csv.NewReader(file)
	citySales := make(map[string]float64)

	rowNum := 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row %d: %w", rowNum, err)
		}
		if len(row) < 3 {
			return nil, fmt.Errorf("row %d is malformed", rowNum)
		}

		// Try to parse the sales column; if it fails on the first row, assume it's a header and skip
		sales, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			if rowNum == 1 {
				rowNum++
				continue // skip header
			}
			return nil, fmt.Errorf("invalid sales value on row %d: %w", rowNum, err)
		}

		city := row[0]
		citySales[city] += sales
		rowNum++
	}

	return &citySales, nil
}

// WriteCitySalesCSV writes the city sales map to the writer in CSV format: City,TotalSales (no header)
func WriteCitySalesCSV(citySales map[string]float64, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	for city, sales := range citySales {
		record := []string{city, strconv.FormatFloat(sales, 'f', -1, 64)}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return writer.Error()
}
