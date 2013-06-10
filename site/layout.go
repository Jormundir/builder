package site

import (
	"builder/site/converter"
	"path/filepath"
)

type layout struct {
	page *page
}

func makeLayout(page *page) *layout {
	return &layout{page: page}
}

func (layout *layout) generateHtml(childContent string, childVars map[string]string, site *Site) (string, error) {
	content, layoutName, vars, err := layout.page.parseSourceFile(layout.page.filepath)
	if err != nil {
		return "", err
	}

	// Turn vars into a map of functions that return the corresponding string value
	tVars := mapToFuncs(vars)
	tVars["page"] = func() map[string]string { return childVars }
	tVars["content"] = func() string { return childContent }

	templatedContent, err := layout.page.templateContent(content, tVars, site.vars())
	if err != nil {
		return "", err
	}
	htmlContent, _, err := converter.ConvertToHtml(templatedContent, filepath.Ext(layout.page.filepath))
	if err != nil {
		return "", err
	}
	parentLayout, err := site.getLayout(layoutName)
	if err != nil {
		return "", err
	}
	if parentLayout != nil {
		return parentLayout.generateHtml(htmlContent, vars, site)
	}
	return htmlContent, nil
}
