package site

import (
	"bufio"
	"builder/site/converter"
	"bytes"
	"errors"
	"html/template"
	"log"
	"os"
	fp "path/filepath"
	"regexp"
	"strings"
)

type page struct {
	fpath string
	fname string
	fext  string

	site   *Site
	vars   template.FuncMap
	layout *layout
	fcont  string

	htmlcont string
	htmlext  string
	page     string
}

func NewPage(fpath string, site *Site) *page {
	page := &page{
		fpath: fpath,
		fname: fp.Base(fpath),
		fext:  fp.Ext(fpath),

		site: site,
		vars: make(template.FuncMap),
	}
	err := page.parseSource()
	if err != nil {
		log.Fatalln(err.Error())
	}
	page.addSmapVar("site", site.Vars)
	return page
}

func (p *page) Build() *page {
	uthtml, ext := converter.Html(p.fcont, p.fext)
	p.addStringVar("path", p.site.webpath(p.fpath, ext))
	err := p.template(uthtml)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if p.layout != nil {
		p.page, p.htmlext, err = p.layout.Render(p.vars, p.htmlcont)
		if err != nil {
			log.Fatalln(err.Error())
		}
		return p
	}
	p.page = p.htmlcont
	p.htmlext = ext
	return p
}

func (p *page) addInterfaceVar(name string, in interface{}) {
	p.vars[name] = in
}

func (p *page) addSmapVar(name string, val map[string]string) {
	p.vars[name] = func() map[string]string { return val }
}

func (p *page) addStringVar(name, val string) {
	p.vars[name] = func() template.HTML { return template.HTML(val) }
}

func (p *page) parseSource() error {
	file, err := os.Open(p.fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := make([]string, 0, 50)
	divided := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		divider, err := regexp.MatchString("^[-]{3,}$", line)
		if err != nil {
			return err
		}
		if divider && !divided {
			divided = true
			err = p.parseVars(lines)
			if err != nil {
				return err
			}
			lines = make([]string, 0, 50)
		} else {
			lines = append(lines, line)
		}
	}
	p.fcont = strings.Join(lines, "\n")
	return nil
}

func (p *page) parseVars(lns []string) error {
	for _, li := range lns {
		if len(li) == 0 {
			continue
		}
		parts := strings.SplitN(li, ":", 2)
		if len(parts) != 2 {
			return errors.New(p.fpath + " ERROR: parsing variables: " + li)
		}
		vname := strings.Trim(parts[0], " ")
		vval := strings.Trim(parts[1], " ")
		if vname == "layout" {
			p.layout = p.site.getLayout(vval)
		} else {
			p.addStringVar(vname, vval)
		}
	}
	return nil
}

func (p *page) template(cont string) error {
	tpl := template.New(p.fpath)
	tpl.Funcs(p.vars)
	tpl, err := tpl.Parse(cont)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, nil)
	if err != nil {
		return err
	}
	p.htmlcont = buf.String()
	return nil
}
