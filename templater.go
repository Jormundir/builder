package builder

import (
	"bytes"
	"html/template"
)

type htmlTemplater struct {
	templated string
}

func newHtmlTemplater() *htmlTemplater {
	return &htmlTemplater{}
}

func (ht *htmlTemplater) template(filepath, content string, pVars, sVars map[string]string) error {
	funcs := make(template.FuncMap)
	if contentVar, ok := pVars[CONTENT]; ok {
		funcs[CONTENT] = func() template.HTML { return template.HTML(contentVar) }
		delete(pVars, CONTENT)
	}
	for name, val := range pVars {
		funcs[name] = func() string { return val }
	}
	funcs[SITE] = func() map[string]string { return sVars }

	tmpl := template.New(filepath)
	tmpl.Funcs(funcs)
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, nil)
	if err != nil {
		return err
	}
	ht.templated = buffer.String()
	return nil
}
