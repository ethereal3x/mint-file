package service

import (
	"github.com/ethereal3x/mint-file/service/parser"
	"path/filepath"
	"strings"
)

type Parser interface {
	Parse(filePath string) (string, error)
	ParseByData(data []byte) (string, error)
}

func GetFileExtension(filePath string) string {
	return strings.ToLower(filepath.Ext(filePath))
}

func GetParserByExtension(ext string) Parser {
	switch ext {
	case ".txt":
		return &parser.TXTParser{}
	case ".md":
		return &parser.MDParser{}
	case ".pdf":
		return &parser.PDFParser{}
	case ".docx":
		return &parser.DocxParser{}
	case ".xlsx":
		return &parser.XLSXParser{}
	default:
		return nil
	}
}
