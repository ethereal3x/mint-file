package parser

import "os"

type TXTParser struct{}

func (t *TXTParser) Parse(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return t.ParseByData(data)
}

func (t *TXTParser) ParseByData(data []byte) (string, error) {
	return string(data), nil
}
