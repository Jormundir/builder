package converter

import (
	"bytes"
	"github.com/knieriem/markdown"
	"strings"
)

func ConvertToHtml(content, ext string) (string, string, error) {
	if len(strings.Trim(content, " ")) == 0 {
		return "", ext, nil
	}

	switch strings.ToLower(ext) {
	case ".md":
		htmlContent := markdownToHtml(content)
		return htmlContent, ".html", nil
	case ".html":
		return content, ext, nil
	default:
		return content, ext, nil
	}
}

func markdownToHtml(content string) string {
	var mparser = markdown.NewParser(nil)
	stringBuffer := strings.NewReader(content)
	bufferReader := new(bytes.Buffer)
	mparser.Markdown(stringBuffer, markdown.ToHTML(bufferReader))
	return bufferReader.String()
}
