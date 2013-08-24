package builder

import (
	"testing"
)

var (
	path    = "template_test_file.md"
	content = "#This is a header\n" +
		"##This is a subheader with a page variable {{layout}}\n\n" +
		"How does a site variable work? {{site.url}}\n"
	pVars = map[string]string{
		"layout": "monkey layout",
	}
	sVars = map[string]string{
		"url": "www.testing.com",
	}

	expectedContent = "#This is a header\n" +
		"##This is a subheader with a page variable " + pVars["layout"] + "\n\n" +
		"How does a site variable work? " + sVars["url"] + "\n"
)

func TestTemplate(t *testing.T) {
	t.Log("called")
	templater := newHtmlTemplater()
	err := templater.template(path, content, pVars, sVars)
	if err != nil {
		t.Fatal(err)
	}
	if templater.templated != expectedContent {
		t.Fatal("Templater templated content\n" + templater.templated + "\n\ndoes not match expected output\n" + expectedContent)
	}
}
