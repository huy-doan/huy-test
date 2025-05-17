package task

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
)

type ReadCSVWithHeaderTask struct{}

func NewReadCSVWithHeaderTask() *ReadCSVWithHeaderTask {
	return &ReadCSVWithHeaderTask{}
}

// Do reads a CSV from io.Reader and returns records as []map[string]string
func (t *ReadCSVWithHeaderTask) Do(r io.Reader) ([]map[string]string, error) {
	var records [][]string
	reader := csv.NewReader(r)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	// Remove BOM if present
	if len(headers) > 0 {
		headers[0] = strings.TrimPrefix(headers[0], "\uFEFF")
	}
	for i, h := range headers {
		headers[i] = strings.TrimSpace(h)
	}
	log.Printf("[DEBUG] CSV headers: %v", headers)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		for len(record) < len(headers) {
			record = append(record, "")
		}
		records = append(records, record)
	}
	// Map headers from Japanese to normalized json (english) form
	mappingTask := NewCSVHeaderMappingTask()
	_, mappedRecords := mappingTask.Do(headers, records)
	return mappedRecords, nil
}
