package builder

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Site struct {
	config *SiteConfig

	parent   *Site
	children map[string]*Site

	pages   map[string]*page
	layouts map[string]*layout

	dirs []string
}

func NewSite(config *SiteConfig, parent *Site) *Site {
	site := &Site{
		config: config,

		parent:   parent,
		children: make(map[string]*Site),

		pages:   make(map[string]*page),
		layouts: make(map[string]*layout),
		dirs:    make([]string, 0, 15),
	}
	return site.init().load()
}

func (s *Site) init() *Site {
	s.children = make(map[string]*Site)
	s.pages = make(map[string]*page)
	s.layouts = make(map[string]*layout)
	s.dirs = make([]string, 0, len(s.dirs))
	err := filepath.Walk(s.config.SourceDir, s.initWalk)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Site) initWalk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	relpath, err := filepath.Rel(s.config.SourceDir, path)
	if err != nil {
		return err
	}

	skippable := s.skippable(relpath)
	if info.IsDir() {
		s.dirs = append(s.dirs, path)
		if path == s.config.SourceDir {
			return nil
		}
		if !skippable {
			site := s.makeChild(relpath)
			s.children[relpath] = site
		}
		return filepath.SkipDir
	}
	if skippable {
		return nil
	}

	page, err := newPage(path, s)
	if err != nil {
		return err
	}
	s.pages[filepath.ToSlash(relpath)] = page

	return nil
}

func (s *Site) load() *Site {
	for path, page := range s.pages {
		_, err := page.build()
		if err != nil {
			log.Fatalln(err)
		}
		delete(s.pages, path)
		s.pages[s.webpath(path, page.webExt)] = page
	}
	return s
}

func (s *Site) GenerateFiles() error {
	os.RemoveAll(s.config.TargetDir)
	for _, child := range s.children {
		err := child.GenerateFiles()
		if err != nil {
			return err
		}
	}
	for webpath, page := range s.pages {
		dest := filepath.Join(s.config.TargetDir, filepath.FromSlash(webpath))
		err := makeDirsTo(dest, DIR_MODE)
		if err != nil {
			log.Println(err)
			continue
		}
		file, err := os.Create(dest)
		if err != nil {
			log.Println(err)
			continue
		}
		defer file.Close()
		_, err = file.WriteString(page.content)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (s *Site) InitGenerate() {
	s.init().load()
	if err := s.GenerateFiles(); err != nil {
		log.Println(err)
	}
}

func (s *Site) MakeWatcher(onEvent func()) *watcher {
	return newWatcher(s, onEvent)
}

func (s *Site) MakeWebserver() *webserver {
	return newWebServer(s)
}

func (s *Site) getDirs() []string {
	dirs := s.dirs
	for _, child := range s.children {
		dirs = append(dirs, child.dirs...)
	}
	return dirs
}

func (s *Site) makeChild(path string) *Site {
	config := s.config.copy()
	config.appendPaths(path)
	child := NewSite(config, s)
	return child
}

func (s *Site) makeLayout(name string) (*layout, error) {
	if layout, ok := s.layouts[name]; ok {
		return layout, nil
	}

	globPath := filepath.Join(s.config.SourceDir, s.config.LayoutDir, name) + ".*"
	matches, err := filepath.Glob(globPath)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		if s.parent != nil {
			return s.parent.makeLayout(name)
		}
		return nil, errors.New("Could not find layout: " + globPath)
	}
	if len(matches) != 1 {
		return nil, errors.New("Ambiguous layout name: " + globPath)
	}
	layout, err := newLayout(matches[0], s)
	if err != nil {
		return nil, err
	}
	s.layouts[name] = layout
	return layout, nil
}

func (s *Site) skippable(path string) bool {
	if path == "." {
		return false
	}
	match, err := superMatch(path, "_*", ".*", s.config.LayoutDir)
	if err != nil {
		log.Fatal(path + " Error checking against ignore path: " + err.Error())
	}
	return match
}

func (s *Site) vars() map[string]string {
	return map[string]string{
		"url": s.config.Url,
	}
}

func (s *Site) webpage(webpath string) *page {
	if webpath == WEB_ROOT {
		if page, ok := s.pages[WEB_DEFAULT_ROOT]; ok {
			return page
		}
		return nil
	}
	parts := strings.Split(webpath, "/")
	if len(parts) == 1 {
		if page, ok := s.pages[parts[0]]; ok {
			return page
		}
	} else if len(parts) > 1 {
		if child, ok := s.children[parts[0]]; ok {
			return child.webpage(strings.Join(parts[1:], "/"))
		}
	}
	return nil
}

func (s *Site) webpath(path, webext string) string {
	webbase := filepath.ToSlash(path)
	webbase = webbase[0 : len(path)-len(filepath.Ext(path))]
	return webbase + webext
}
