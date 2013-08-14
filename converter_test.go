package builder

import (
	"path/filepath"
	"testing"
)

var (
	markdownExample = "#This is example markdown\n" +
		"Just a paragraph\n\n" +
		"## subheader"
	markdownExt             = ".md"
	markdownExampleExpected = "<h1>This is example markdown</h1>\n\n" +
		"<p>Just a paragraph</p>\n\n" +
		"<h2>subheader</h2>\n"
	markdownExtExpected = ".html"

	cssExample = "* {\n" +
		"\tmargin: 0;\n" +
		"\tpadding: 0;\n"
	cssExt             = ".css"
	cssExampleExpected = "* {\n" +
		"\tmargin: 0;\n" +
		"\tpadding: 0;\n"
	cssExtExpected = ".css"

	filepathExample         = filepath.FromSlash("example/file/path")
	filepathExampleWebExt   = ".html"
	filepathExampleExpected = "example/file/path.html"
)

func TestToHtml(t *testing.T) {
	converter := newWebConverter()

	converter.toHtml(markdownExample, markdownExt)
	if converter.webBody != markdownExampleExpected {
		t.Error("Markdown conversion did not produce expected output.")
		t.Fail()
	}
	if converter.webExt != markdownExtExpected {
		t.Error("Markdown conversion did not produce expected extension")
		t.Fail()
	}

	converter.toHtml(cssExample, cssExt)
	if converter.webBody != cssExampleExpected {
		t.Error("CSS conversion did not produce expected output")
		t.Fail()
	}
	if converter.webExt != cssExtExpected {
		t.Error("CSS conversion did not produce expected extension")
		t.Fail()
	}
}

func TestWebpath(t *testing.T) {
	converter := newWebConverter()

	webpathActual := converter.webpath(filepathExample, filepathExampleWebExt)
	if webpathActual != filepathExampleExpected {
		t.Error("Filepath to Webpath conversion did not produce expected result")
		t.Fail()
	}
}
