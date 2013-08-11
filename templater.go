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
	for name, val := range pVars {
		if name == CONTENT {
			funcs[name] = func() template.HTML { return template.HTML(val) }
		} else {
			funcs[name] = func() string { return val }
		}
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
