package builder

type layout struct {
	*sourceParser
	*webConverter
	*htmlTemplater

	parent *layout
	site   *Site
}

func newLayout(path string, site *Site) (*layout, error) {
	l := &layout{
		sourceParser:  newSourceParser(path),
		webConverter:  newWebConverter(),
		htmlTemplater: newHtmlTemplater(),
		site:          site,
	}
	return l.init()
}

func (l *layout) init() (*layout, error) {
	if err := l.parse(); err != nil {
		return nil, err
	}
	if layout, ok := l.vars[LAYOUT]; ok {
		delete(l.vars, LAYOUT)
		parent, err := l.site.makeLayout(layout)
		if err != nil {
			return nil, err
		}
		l.parent = parent
	}
	return l, nil
}

func (l *layout) render(pVars map[string]string, content string) (string, error) {
	l.vars[CONTENT] = content
	l.toHtml(l.body, l.ext)
	if err := l.template(l.path, l.webBody, l.vars, l.site.vars()); err != nil {
		return "", err
	}
	if l.parent != nil {
		return l.parent.render(pVars, l.templated)
	}
	return l.templated, nil
}
