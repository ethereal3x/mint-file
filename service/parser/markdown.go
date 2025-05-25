package parser

import "os"

type MDParser struct{}

// Parse 通过文件路径读取 Markdown 文件内容
func (m *MDParser) Parse(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return m.ParseByData(data)
}

// ParseByData 从内存数据（[]byte）读取 Markdown 内容
func (m *MDParser) ParseByData(data []byte) (string, error) {
	return string(data), nil
}
