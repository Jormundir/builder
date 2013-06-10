package site

import (
	"bufio"
	"builder/site/converter"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type page struct {
	filepath        string
	fileContent     string
	layout          *layout
	vars            map[string]string
	baseHtmlContent string
	fullHtmlContent string
}

func makePage(path string, site *Site) (*page, error) {
	page := &page{filepath: path}
	fileContent, layoutName, vars, err := page.parseSourceFile(page.filepath)
	if err != nil {
		return nil, err
	}
	page.fileContent = fileContent
	page.layout, err = site.getLayout(layoutName)
	if err != nil {
		return nil, err
	}
	page.vars = vars

	untemplatedHtmlContent, err := converter.ConvertToHtml(page.fileContent, filepath.Ext(path))
	if err != nil {
		return nil, err
	}
	baseHtmlContent, err := page.templateContent(untemplatedHtmlContent, mapToFuncs(page.vars), site.vars())
	if err != nil {
		return nil, err
	}
	page.baseHtmlContent = baseHtmlContent

	var fullHtmlContent string
	if page.layout != nil {
		fullHtmlContent, err = page.layout.generateHtml(page.baseHtmlContent, page.vars, site)
		if err != nil {
			return nil, err
		}
	} else {
		fullHtmlContent = page.baseHtmlContent
	}
	page.fullHtmlContent = fullHtmlContent
	return page, nil
}

func (page *page) parseSourceFile(filepath string) (fileContent, layoutName string, vars map[string]string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	fileLines := make([]string, 0, 50)
	encounteredVarsDivider := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var lineIsVarsDivider bool
		lineIsVarsDivider, err = regexp.MatchString("^[-]{3,}$", line)
		if err != nil {
			return
		}

		if lineIsVarsDivider && !encounteredVarsDivider {
			encounteredVarsDivider = true
			vars, err = page.parseVars(filepath, fileLines)
			if err != nil {
				return
			}
			fileLines = make([]string, 0, 50)
		} else {
			// Trim initial empty lines
			if len(strings.Trim(line, " ")) == 0 {
				continue
			}
			fileLines = append(fileLines, line)
		}
	}

	var layoutSpecified bool
	layoutName, layoutSpecified = vars["layout"]
	if layoutSpecified {
		delete(vars, "layout")
	}
	fileContent = strings.Join(fileLines, "\n")
	return
}

func (page *page) parseVars(path string, varLines []string) (vars map[string]string, err error) {
	vars = make(map[string]string)
	for _, line := range varLines {
		// skip empty lines, they don't have a variable declaration
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			err = variableDeclarationError{op: "Parsing Variables", path: path}
			return
		}

		varName := strings.Trim(parts[0], " ")
		varValue := strings.Trim(parts[1], " ")
		vars[varName] = varValue
	}
	return
}

func (page *page) templateContent(content string, vars map[string]interface{}, siteVars map[string]string) (string, error) {
	var err error
	tplt := template.New("template")
	vars["site"] = func() map[string]string { return siteVars }
	tplt.Funcs(vars)
	tplt, err = tplt.Parse(content)
	if err != nil {
		return "", err
	}

	stringHolder := new(stringHolder)
	err = tplt.Execute(stringHolder, nil)
	if err != nil {
		return "", err
	}
	content = stringHolder.string()
	return content, nil
}

// STRINGHOLDER SEEMS REALLY STUPID...
type stringHolder struct {
	str []byte
}

func (str *stringHolder) string() string {
	return string(str.str)
}

func (str *stringHolder) Write(p []byte) (int, error) {
	str.str = append(str.str, p...)
	return len(str.str), nil
}
