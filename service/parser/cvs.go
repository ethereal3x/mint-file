package parser

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
)

type CSVParser struct{}

func (c *CSVParser) Parse(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for _, row := range records {
		builder.WriteString(strings.Join(row, ", ") + "\n")
	}
	return builder.String(), nil
}

func (c *CSVParser) ParseByData(data []byte) (string, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for _, row := range records {
		builder.WriteString(strings.Join(row, ", ") + "\n")
	}
	return builder.String(), nil
}
