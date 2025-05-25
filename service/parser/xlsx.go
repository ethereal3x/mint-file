package parser

import (
	"bytes"
	"github.com/xuri/excelize/v2"
	"os"
	"strings"
)

type XLSXParser struct{}

// Parse 从文件路径解析 .xlsx 文件内容
func (x *XLSXParser) Parse(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return x.ParseByData(data)
}

func (x *XLSXParser) ParseByData(data []byte) (string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return "", err
		}
		for _, row := range rows {
			builder.WriteString(strings.Join(row, "\t") + "\n")
		}
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
