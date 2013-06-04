package main

import (
	"bufio"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

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

type siteFile struct {
	vars        map[string]string
	content     string
	fullContent string
}

func NewSiteFile(path string, siteBuilder *SiteBuilder, siteVars SiteVars) (*siteFile, error) {
	// pars vars and content
	vars, content, err := parse(path)
	if err != nil {
		return nil, err
	}

	// take template out of vars...
	templ, templateSpecified := vars["template"]
	if templateSpecified {
		delete(vars, "template")
	}
	templatePath := siteBuilder.templatePath(templ)

	// convert content to html
	content, err = convertToHTML(content, filepath.Ext(path))
	if err != nil {
		return nil, err
	}

	// Add site vars to page vars for template access
	//for key, val := range siteVars {
	//	vars["site."+key] = val
	//}

	// run content through templater.
	tplt := template.New("template")
	tplt.Funcs(template.FuncMap{"site": func() SiteVars { return siteVars }})
	tplt, err = tplt.Parse(content)
	if err != nil {
		return nil, err
	}

	stringHolder := new(stringHolder)
	err = tplt.Execute(stringHolder, vars)
	if err != nil {
		return nil, err
	}
	content = stringHolder.string()

	// if a template was specified, pass the content to the template to make the fullContent
	var fullContent string
	if templateSpecified {
		htmlContent := template.HTML(content)
		fullContent, err = renderTemplate(templatePath, htmlContent, vars, siteBuilder, siteVars)
		if err != nil {
			return nil, err
		}
	} else {
		fullContent = content
	}

	return &siteFile{vars: vars, content: content, fullContent: fullContent}, nil
}

func (file *siteFile) getContent() string {
	return file.content
}

// yes yes, this is almost an exact double of NewSiteFile...
func renderTemplate(path string, childContent template.HTML, childVars map[string]string, siteBuilder *SiteBuilder, siteVars SiteVars) (string, error) {
	tVars, tContent, err := parse(path)
	if err != nil {
		return "", err
	}

	// take template out of vars if it was declared
	templ, templateSpecified := tVars["template"]
	if templateSpecified {
		delete(tVars, "template")
	}
	templatePath := siteBuilder.templatePath(templ)

	// make variables for template engine..
	var vars struct {
		Child    map[string]string
		Content  template.HTML
		Template map[string]string
	}
	vars.Child = childVars
	vars.Template = tVars
	vars.Content = childContent

	// run content through templater.
	tplt := template.New("template")
	tplt.Funcs(template.FuncMap{"site": func() SiteVars { return siteVars }})
	tplt, err = tplt.Parse(tContent)
	if err != nil {
		return "", err
	}
	stringHolder := new(stringHolder)
	err = tplt.Execute(stringHolder, vars)
	if err != nil {
		return "", err
	}
	content := stringHolder.string()

	// convert content to html
	content, err = convertToHTML(content, filepath.Ext(path))
	if err != nil {
		return "", err
	}

	if templateSpecified {
		// prep vars for parent template...combine child and template's vars into one map - may need more thinking.
		for key, val := range tVars {
			childVars[key] = val
		}
		htmlContent := template.HTML(content)
		return renderTemplate(templatePath, htmlContent, childVars, siteBuilder, siteVars)
	} else {
		return content, nil
	}
}

func parse(path string) (vars map[string]string, content string, err error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return vars, content, err
	}
	defer file.Close()

	// parse vars and content
	fileLines := make([]string, 0, 50)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineIsVarDivider, err := regexp.MatchString("^[-]{3,}$", line)
		if err != nil {
			return vars, content, err
		}

		if lineIsVarDivider {
			vars, err = parseVars(path, fileLines)
			if err != nil {
				return vars, content, err
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

	content = strings.Join(fileLines, "\n")
	return vars, content, err
}

func parseVars(path string, varLines []string) (map[string]string, error) {
	vars := make(map[string]string)
	var err error
	for num, line := range varLines {
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			err = errors.New("Invalid variable declaration syntax in " + path + " line " + string(num+1))
			return vars, err
		}

		varName := strings.Trim(parts[0], " ")
		varValue := strings.Trim(parts[1], " ")
		vars[varName] = varValue
	}

	return vars, err
}
