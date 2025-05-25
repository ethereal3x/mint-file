package parser

import (
	"bytes"
	"github.com/ledongthuc/pdf"
	"os"
)

type PDFParser struct{}

func (p *PDFParser) Parse(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var text string
	totalPage := r.NumPage()
	for i := 1; i <= totalPage; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		text += content
	}
	return text, nil
}

func (p *PDFParser) ParseByData(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	r, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", err
	}

	var text string
	totalPage := r.NumPage()
	for i := 1; i <= totalPage; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		text += content
	}
	return text, nil
}
