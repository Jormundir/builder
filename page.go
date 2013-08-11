package builder

type page struct {
	*sourceParser
	*webConverter
	*htmlTemplater

	site    *Site
	layout  *layout
	content string
}

func newPage(path string, site *Site) (*page, error) {
	p := &page{
		sourceParser:  newSourceParser(path),
		webConverter:  newWebConverter(),
		htmlTemplater: newHtmlTemplater(),
		site:          site,
	}
	return p.init()
}

func (p *page) init() (*page, error) {
	err := p.parse()
	if err != nil {
		return nil, err
	}
	if layout, ok := p.vars[LAYOUT]; ok {
		delete(p.vars, LAYOUT)
		p.layout, err = p.site.makeLayout(layout)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *page) build() (*page, error) {
	p.toHtml(p.body, p.ext)
	err := p.template(p.path, p.webBody, p.vars, p.site.vars())
	if err != nil {
		return nil, err
	}
	if p.layout != nil {
		p.content, err = p.layout.render(p.vars, p.templated)
		if err != nil {
			return nil, err
		}
	} else {
		p.content = p.templated
	}
	return p, nil
}
