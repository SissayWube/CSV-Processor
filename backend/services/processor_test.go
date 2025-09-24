package services

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestProcessCSV(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		csvData     string
		expected    *map[string]float64
		expectedErr bool
	}{
		{
			name: "Valid CSV Data",
			csvData: `City,Product,Sales
New York,Laptop,1200.50
Los Angeles,Tablet,800.00
New York,Tablet,750.00`,
			expected:    &map[string]float64{"New York": 1950.50, "Los Angeles": 800.00},
			expectedErr: false,
		},
		{
			name: "Invalid Sales Value",
			csvData: `City,Product,Sales
New York,Laptop,abc
Los Angeles,Tablet,800.00`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name: "Malformed Row",
			csvData: `City,Product
New York,Laptop,1200.50`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Empty CSV Data",
			csvData:     ``,
			expected:    &map[string]float64{},
			expectedErr: false,
		},
		{
			name: "CSV with header only",
			csvData: `City,Product,Sales`,
			expected:    &map[string]float64{},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a reader from the CSV data
			reader := strings.NewReader(tc.csvData)

			// Process the CSV data
			result, err := ProcessCSV(reader)

			// Check for error
			if tc.expectedErr {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Compare the result with the expected value
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Result mismatch: got %v, expected %v", result, tc.expected)
			}
		})
	}
}

func TestWriteCitySalesCSV(t *testing.T) {
	// Test case
	// NOTE: The order of lines in expectedCSV does not matter because we sort before comparing.
	citySales := map[string]float64{"New York": 1950.50, "Los Angeles": 800.00}
	expectedCSV := "New York,1950.5\nLos Angeles,800\n"

	// Create a buffer to write the CSV data to
	buffer := new(bytes.Buffer)

	// Write the city sales data to the buffer
	err := WriteCitySalesCSV(citySales, buffer)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Get the CSV data from the buffer
	csvData := buffer.String()

	// Split the CSV data into lines and sort them
	expectedLines := strings.Split(strings.TrimSpace(expectedCSV), "\n")
	actualLines := strings.Split(strings.TrimSpace(csvData), "\n")

	// Sort the lines for comparison
	// Sort the lines for comparison (since map order is not deterministic)
	if len(actualLines) != len(expectedLines) {
		t.Fatalf("Result mismatch: number of lines differ. Got %d, expected %d", len(actualLines), len(expectedLines))
	}
	// Sort the lines for comparison (since map order is not deterministic).
	sort.Strings(expectedLines)
	sort.Strings(actualLines)
	if !reflect.DeepEqual(actualLines, expectedLines) {
		t.Fatalf("Result mismatch:\ngot:\n%s\nexpected:\n%s", strings.Join(actualLines, "\n"), strings.Join(expectedLines, "\n"))
	}

	// Test case with empty map
	citySalesEmpty := map[string]float64{}
	expectedCSVEmpty := ""

	bufferEmpty := new(bytes.Buffer)
	errEmpty := WriteCitySalesCSV(citySalesEmpty, bufferEmpty)
	if errEmpty != nil {
		t.Fatalf("Unexpected error: %v", errEmpty)
	}

	csvDataEmpty := bufferEmpty.String()
	if csvDataEmpty != expectedCSVEmpty {
		t.Fatalf("Result mismatch for empty map: got %s, expected %s", csvDataEmpty, expectedCSVEmpty)
	}
}