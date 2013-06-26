package site

import (
	"log"
	"os"
	"path/filepath"
)

const DIR_MODE = 0777

type Site struct {
	config  *siteConfig
	dirs    []string
	layouts map[string]*layout
	pages   map[string]*page

	Vars map[string]string
}

func NewSite(cmdln map[string]*string) *Site {
	config := NewConfig(cmdln)
	s := &Site{
		config:  config,
		dirs:    make([]string, 0, 10),
		layouts: make(map[string]*layout),
		pages:   make(map[string]*page),
		Vars:    make(map[string]string),
	}
	s.Vars["url"] = s.config.SiteUrl
	return s.Init()
}

func (s *Site) Execute() {
	s.Build()
	s.webMode()
	s.Init()

	webserver := newWebserver(s)
	webserver.serve()
	log.Println("helo")
	for {
	}

	//watcher := newWatcher(s)
	//watcher.watch()
}

// Would be cool to make a backup, and if anything breaks during build,
// remake the target dir into what it was before the build.
// Backup should happen in execute so the watcher doesn't backup
// with every file change
func (s *Site) Build() *Site {
	os.RemoveAll(s.config.TargetDir)
	for path, page := range s.pages {
		dest := filepath.Join(s.config.TargetDir, path)
		err := makeDirsTo(dest, DIR_MODE)
		if err != nil {
			log.Print(err)
			continue
		}
		pf, err := os.Create(dest)
		if err != nil {
			log.Print(err)
			continue
		}
		defer pf.Close()
		_, err = pf.WriteString(page.page)
		if err != nil {
			log.Print(err)
			continue
		}
	}
	return s
}

func (s *Site) Init() *Site {
	filepath.Walk(s.config.SourceDir,
		func(p string, i os.FileInfo, _ error) error {
			ignr := s.ignorePath(p)
			if i.IsDir() {
				s.dirs = append(s.dirs, p)
				if ignr {
					return filepath.SkipDir
				}
				return nil
			}
			if ignr {
				return nil
			}
			page := NewPage(p, s).Build() // worry about aggregating errors later.
			prel := s.relPath(page.fpath, page.htmlext)
			s.pages[prel] = page
			return nil
		})
	return s
}

func (s *Site) getLayout(nm string) *layout {
	if len(nm) == 0 {
		log.Println("layout " + nm + " cannot be found.")
		return nil
	}
	lyt, ok := s.layouts[nm]
	if ok {
		return lyt
	}
	gp := filepath.Join(s.config.SourceDir, s.config.LayoutDir, nm) + ".*"
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
	lyt = newLayout(matches[0], s)
	s.layouts[nm] = lyt
	return lyt
}

func (s *Site) ignorePath(p string) bool {
	rp, err := filepath.Rel(s.config.SourceDir, p)
	if err != nil {
		log.Fatal(p + " Error checking ignore relative path")
	}
	if rp == "." {
		return false
	}
	match, err := superMatch(rp, "_*", ".*", s.config.LayoutDir)
	if err != nil {
		log.Fatal(p + " Error checking against ignore path")
	}
	return match
}

func (s *Site) relPath(p, ext string) string {
	rp, err := filepath.Rel(s.config.SourceDir, p)
	if err != nil {
		log.Fatal(p + " Error getting relative path")
	}
	extLen := len(filepath.Ext(p))
	rp = rp[0:(len(rp) - extLen)]
	return rp + ext
}

func (s *Site) webpath(p, ext string) string {
	return "/" + filepath.ToSlash(s.relPath(p, ext))
}

func (s *Site) buildMode() {
	s.Vars["url"] = s.config.SiteUrl
}

func (s *Site) webMode() {
	s.Vars["url"] = "http://localhost:1337"
}
