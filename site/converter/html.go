package converter

import (
	"bytes"
	"github.com/knieriem/markdown"
	"strings"
)

func Html(cont, ext string) (string, string) {
	if len(strings.Trim(cont, " ")) == 0 {
		return "", ext
	}

	switch strings.ToLower(ext) {
	case ".md":
		htmlcont := markdownToHtml(cont)
		return htmlcont, ".html"
	case ".html":
		return cont, ext
	default:
		return cont, ext
	}
}

func markdownToHtml(cont string) string {
	var mparser = markdown.NewParser(nil)
	stringBuffer := strings.NewReader(cont)
	bufferReader := new(bytes.Buffer)
	mparser.Markdown(stringBuffer, markdown.ToHTML(bufferReader))
	return bufferReader.String()
}
