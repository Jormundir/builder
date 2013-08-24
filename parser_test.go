package builder

import (
	"os"
	"testing"
)

var (
	filename    = "parser_test_file.md"
	ext         = ".md"
	testContent = "layout: hello\n" +
		"width: 10px\n" +
		"---\n" +
		"## this is some page content\n" +
		"more page content\n" +
		"{{width}} variable use\n"
	testExpectedVars = map[string]string{
		"layout": "hello",
		"width":  "10px",
	}
	testExpectedBody = "## this is some page content\n" +
		"more page content\n" +
		"{{width}} variable use"
)

func TestParse(t *testing.T) {
	// make file
	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		file.Close()
		os.Remove(filename)
	}()
	// fill with test contents
	_, err = file.WriteString(testContent)
	if err != nil {
		t.Fatal(err)
	}
	// parse test file
	parser := newSourceParser(filename)
	err = parser.parse()
	if err != nil {
		t.Fatal(err)
	}
	// make sure variables and content are what's expected
	switch {
	case parser.path != filename:
		t.Fatal("Parser path " + parser.path + " does not match input file path " + filename)
	case parser.ext != ext:
		t.Fatal("Parser extension " + parser.ext + " does not match expected " + ext)
	case !mapsEqual(parser.vars, testExpectedVars):
		t.Fatalf("Parser vars %v do not match expected vars %v", parser.vars, testExpectedVars)
	case parser.body != testExpectedBody:
		t.Fatal("Parser body \n" + parser.body + "\n\ndoes not match expected\n" + testExpectedBody)
	}
}

func mapsEqual(map1, map2 map[string]string) bool {
	for index, value := range map1 {
		corresponding, ok := map2[index]
		if !ok || value != corresponding {
			return false
		}
	}
	return true
}
