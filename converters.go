package main

import (
	"bytes"
	"github.com/knieriem/markdown"
	"strings"
)

func convertToHTML(content, ext string) (string, error) {
	if len(strings.Trim(content, " ")) == 0 {
		return "", nil
	}

	switch strings.ToLower(ext) {
	case ".md":
		return markdownToHTML(content)
	case ".html":
		return content, nil
	default:
		return content, nil
	}
}

func markdownToHTML(content string) (string, error) {
	var mparser = markdown.NewParser(nil)
	stringBuffer := strings.NewReader(content)
	bufferReader := new(bytes.Buffer)
	mparser.Markdown(stringBuffer, markdown.ToHTML(bufferReader))
	return bufferReader.String(), nil
}
