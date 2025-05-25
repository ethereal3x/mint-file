package parser

import (
	"bytes"
	"encoding/xml"
	"github.com/nguyenthenguyen/docx"
	"io"
	"strings"
)

type DocxParser struct{}

func extractTextByXML(data []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var sb strings.Builder

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// 找到 <w:t> 标签的开始元素
		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "t" {
				// 读取 <w:t> 标签内的字符数据
				var content string
				if err := decoder.DecodeElement(&content, &se); err != nil {
					return "", err
				}
				sb.WriteString(content)
			}
		}
	}

	return sb.String(), nil
}

func (w *DocxParser) Parse(filePath string) (string, error) {
	doc, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", err
	}
	defer doc.Close()

	content := doc.Editable().GetContent()

	text, err := extractTextByXML([]byte(content))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

func (w *DocxParser) ParseByData(data []byte) (string, error) {
	r, err := docx.ReadDocxFromMemory(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	defer r.Close()

	content := r.Editable().GetContent()

	text, err := extractTextByXML([]byte(content))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}
