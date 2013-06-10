package converter

import (
	"bytes"
	"github.com/knieriem/markdown"
	"strings"
)

func ConvertToHtml(content, ext string) (string, error) {
	if len(strings.Trim(content, " ")) == 0 {
		return "", nil
	}

	switch strings.ToLower(ext) {
	case ".md":
		return markdownToHtml(content)
	case ".html":
		return content, nil
	default:
		return content, nil
	}
}

func markdownToHtml(content string) (string, error) {
	var mparser = markdown.NewParser(nil)
	stringBuffer := strings.NewReader(content)
	bufferReader := new(bytes.Buffer)
	mparser.Markdown(stringBuffer, markdown.ToHTML(bufferReader))
	return bufferReader.String(), nil
}
