package builder

import (
	"bytes"
	"github.com/knieriem/markdown"
	"path/filepath"
	"strings"
)

type webConverter struct {
	webBody string
	webExt  string
}

func newWebConverter() *webConverter {
	return &webConverter{}
}

func (wc *webConverter) toHtml(content, ext string) {
	switch strings.ToLower(ext) {
	case ".md":
		html := wc.markdownToHtml(content)
		wc.webBody = html
		wc.webExt = ".html"
	default:
		wc.webBody = content
		wc.webExt = ext
	}
}

func (wc *webConverter) markdownToHtml(content string) string {
	var mparser = markdown.NewParser(nil)
	stringBuffer := strings.NewReader(content)
	bufferReader := new(bytes.Buffer)
	mparser.Markdown(stringBuffer, markdown.ToHTML(bufferReader))
	return bufferReader.String()
}

func (wc *webConverter) webpath(fpath, webext string) string {
	webbase := filepath.ToSlash(fpath)
	webbase = webbase[0 : len(fpath)-len(filepath.Ext(fpath))]
	return webbase + webext
}
