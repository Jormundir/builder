package site

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Site struct {
	config  *siteConfig
	dirs    []string
	layouts map[string]*layout
	pages   map[string]*page

	Vars map[string]string
}

func NewSite(cmdln map[string]*string) *Site {
	config := NewConfig(cmdln)
	site := &Site{
		config:  config,
		dirs:    make([]string, 0, 10),
		layouts: make(map[string]*layout),
		pages:   make(map[string]*page),
		Vars:    make(map[string]string),
	}
	return site.Init()
}

func (site *Site) Execute() {
	fmt.Printf("%v\n", site)
}

func (site *Site) Init() *Site {
	filepath.Walk(site.config.SourceDir,
		func(p string, i os.FileInfo, _ error) error {
			ignr := site.ignorePath(p)
			if i.IsDir() {
				site.dirs = append(site.dirs, p)
				if ignr {
					return filepath.SkipDir
				}
				return nil
			}
			if ignr {
				return nil
			}
			page := NewPage(p, site).Build() // worry about aggregating errors later.
			prel := site.relPath(p)
			site.pages[prel] = page
			return nil
		})
	site.Vars["url"] = site.config.SiteUrl
	return site
}

func (site *Site) getLayout(nm string) *layout {
	if len(nm) == 0 {
		log.Println("layout " + nm + " cannot be found.")
		return nil
	}
	lyt, ok := site.layouts[nm]
	if ok {
		return lyt
	}
	gp := filepath.Join(site.config.SourceDir, site.config.LayoutDir, nm) + ".*"
	matches, err := filepath.Glob(gp)
	if err != nil {
		log.Println("Error finding layout files for " + nm)
		return nil
	}
	if len(matches) == 0 {
		log.Println("Layout could not be found: " + nm)
		return nil
	}
	if len(matches) != 1 {
		log.Println("Ambiguous layout name " + nm)
		return nil
	}
	lyt = newLayout(matches[0], site)
	site.layouts[nm] = lyt
	return lyt
}

func (site *Site) ignorePath(p string) bool {
	rp, err := filepath.Rel(site.config.SourceDir, p)
	if err != nil {
		log.Fatal(p + " Error checking ignore relative path")
	}
	if rp == "." {
		return false
	}
	match, err := filepath.Match("_*", p)
	if err != nil {
		log.Fatal(p + " Error checking against ignore path")
	}
	if !match {
		match, err = filepath.Match(".*", p)
		if err != nil {
			log.Fatal(p + " Error checking against ignore path")
		}
	}
	return match
}

func (site *Site) relPath(p string) string {
	rp, err := filepath.Rel(site.config.SourceDir, p)
	if err != nil {
		log.Fatal(p + " Error getting relative path")
	}
	return rp
}
