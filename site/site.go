package site

import (
	"os"
	"path/filepath"
	"strings"
)

const DIR_MODE = 0777

type Site struct {
	config      *siteConfig
	directories []string
	layouts     map[string]*layout
	pages       map[string]*page
}

func NewSite() (*Site, error) {
	// Get configuration
	siteConfig, err := newSiteConfig()
	if err != nil {
		return nil, err
	}
	site := &Site{
		config:      siteConfig,
		directories: make([]string, 0, 10),
		layouts:     make(map[string]*layout),
		pages:       make(map[string]*page),
	}
	err = site.parseSource()
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (site *Site) BuildSite() (errs []error) {
	errs = make([]error, 0, 10)
	// backup target directory
	if site.config.Backup {
		// clear backup directory because path collisions can cause errors. Make more robust later..
		err := site.clearDir(site.config.BackupDir)
		if err != nil {
			return append(errs, err)
		}
		err = site.backupTargetDir()
		if err != nil {
			return append(errs, err)
		}
	}

	// wipe target directory
	err := site.clearDir(site.config.TargetDir)
	if err != nil {
		return append(errs, err)
	}

	// Create page files
	for path, page := range site.pages {
		path = filepath.Join(site.config.TargetDir, path+page.ext)
		MakeDirectoriesTo(path, DIR_MODE)
		pagefile, err := os.Create(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		defer pagefile.Close()
		_, err = pagefile.WriteString(page.fullHtmlContent)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

func (site *Site) backupTargetDir() error {
	return filepath.Walk(site.config.TargetDir, func(path string, info os.FileInfo, _ error) error {
		if info == nil {
			return nil
		}

		backupPath, err := filepath.Rel(site.config.TargetDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			backupDirPath := filepath.Join(site.config.BackupDir, backupPath)
			err := os.MkdirAll(backupDirPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			_, err := CopyFile(filepath.Join(site.config.BackupDir, backupPath), path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (site *Site) clearDir(path string) error {
	return os.RemoveAll(path)
}

func (site *Site) getLayout(name string) (*layout, error) {
	if len(name) == 0 {
		return nil, nil
	}
	layout, ok := site.layouts[name]
	if ok {
		return layout, nil
	}

	layoutGlobPath := filepath.Join(site.config.SourceDir, site.config.LayoutDir, name) + ".*"
	matches, err := filepath.Glob(layoutGlobPath)
	if err != nil {
		return nil, err
	}
	if len(matches) != 1 {
		return nil, ambiguousLayoutNameError{name}
	}
	page, err := makePage(matches[0], site)
	if err != nil {
		return nil, err
	}
	layout = makeLayout(page)
	site.layouts[name] = layout
	return layout, nil
}

// Janky function...
func (site *Site) ignorePath(path string) bool {
	relativePath, err := filepath.Rel(site.config.SourceDir, path)
	if err != nil {
		panic(err.Error())
	}
	if relativePath == "." {
		return false
	}
	slashRelativePath := strings.TrimPrefix(filepath.ToSlash(relativePath), "/")
	pathParts := strings.Split(slashRelativePath, "/")
	ignorePatterns := []string{"_*", ".*", site.config.LayoutDir}
	ignore := false
	for _, part := range pathParts {
		for _, pattern := range ignorePatterns {
			match, err := filepath.Match(pattern, part)
			if err != nil {
				panic(err.Error())
			}
			if match {
				ignore = true
			}
		}
	}
	return ignore
}

func (site *Site) MakeWebServer() *WebServer {
	return &WebServer{pages: site.pages, port: "8080"}
}

func (site *Site) parseSource() error {
	return filepath.Walk(site.config.SourceDir, func(path string, info os.FileInfo, _ error) error {
		ignore := site.ignorePath(path)
		if info.IsDir() {
			site.directories = append(site.directories, path)
			if ignore {
				return filepath.SkipDir
			}
			return nil
		}
		if ignore {
			return nil
		}

		page, err := makePage(path, site)
		if err != nil {
			return err
		}
		relPath, err := site.relPath(site.config.SourceDir, path)
		if err != nil {
			return err
		}
		site.pages[relPath] = page
		return nil
	})
}

func (site *Site) relPath(base, path string) (string, error) {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(rel, filepath.Ext(path)), nil
}

func (site *Site) vars() map[string]string {
	return map[string]string{
		"url": site.config.SiteUrl,
	}
}
