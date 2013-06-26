package site

import (
	tpl "html/template"
	"reflect"
)

type layout struct {
	pg *page
}

func newLayout(p string, s *Site) *layout {
	return &layout{pg: NewPage(p, s)}
}

func (l *layout) Render(vars tpl.FuncMap, cont string) (string, string, error) {
	err := l.pg.parseSource()
	if err != nil {
		return "", "", err
	}
	l.pg.addStringVar("content", cont)
	// hacky hacky...
	page := make(map[string]tpl.HTML)
	for name, fun := range vars {
		if name == "site" {
			continue
		}
		switch n := reflect.TypeOf(fun).Out(0).Name(); n {
		case "HTML":
			var v []reflect.Value
			page[name] = tpl.HTML(reflect.ValueOf(fun).Call(v)[0].Convert(reflect.TypeOf("")).String())
		default:
			continue
		}
	}
	// not proud of this.
	l.pg.vars["page"] = func() map[string]tpl.HTML { return page }
	l.pg.Build()
	return l.pg.page, l.pg.htmlext, nil
}
