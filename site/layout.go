package site

import (
	tpl "html/template"
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
	l.pg.addInterfaceVar("page", vars)
	l.pg.Build()
	return l.pg.page, l.pg.htmlext, nil
}
