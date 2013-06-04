package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DIRECTORY_PERMISSIONS = 0777

type siteConfig struct {
	Backup      bool
	BackupDir   string
	SiteUrl     string
	SourceDir   string
	TargetDir   string
	TemplateDir string
}

func newSiteConfig() *siteConfig {
	siteConfig := &siteConfig{
		Backup:      true,
		BackupDir:   "_backup",
		SiteUrl:     "http://localhost:4000",
		SourceDir:   ".",
		TargetDir:   "_site",
		TemplateDir: "templates",
	}
	siteConfig.parseConfigFile()
	siteConfig.parseCommandLine()
	return siteConfig
}

func (siteConfig *siteConfig) parseConfigFile() {
	configFile, err := os.Open("_config.json")
	if err != nil {
		fmt.Println("Warning: _config.json could not be parse for settings.")
		return
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(siteConfig)
	return
}

func (siteConfig *siteConfig) parseCommandLine() {
	flag.BoolVar(&siteConfig.Backup, "backup", siteConfig.Backup,
		"files in target directory will be backed up to _backup/"+
			" in your source directory.")
	flag.StringVar(&siteConfig.SourceDir, "source", siteConfig.SourceDir,
		"where your project source files are.")
	flag.StringVar(&siteConfig.TargetDir, "target", siteConfig.TargetDir,
		"where to build your site")
	flag.Parse()
}

type SiteVars map[string]string

type SiteBuilder struct {
	config *siteConfig
	vars   *SiteVars
}

func NewSiteBuilder() *SiteBuilder {
	siteBuilder := &SiteBuilder{
		config: newSiteConfig(),
	}
	siteBuilder.vars = &SiteVars{} // TODO: Figure out site variables.
	return siteBuilder
}

func (siteBuilder *SiteBuilder) buildSite() {
	// backup target directory
	if siteBuilder.config.Backup {
		err := siteBuilder.backupTargetDir()
		if err != nil {
			fmt.Println("Error backing up target directory: " + err.Error())
		}
	}

	// wipe target directory
	err := siteBuilder.clearTargetDir()
	if err != nil {
		fmt.Println("Error clearing target directory: " + err.Error())
		return
	}

	// crawl source directory
	// Need a better way to do this
	err = siteBuilder.buildSiteFiles()
	if err != nil {
		fmt.Println("Error building site files: " + err.Error())
		return
	}
}

func (siteBuilder *SiteBuilder) backupTargetDir() error {
	return filepath.Walk(siteBuilder.config.TargetDir, func(path string, info os.FileInfo, _ error) error {
		if info == nil {
			return nil
		}

		backupPath, err := filepath.Rel(siteBuilder.config.TargetDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			backupDirPath := filepath.Join(siteBuilder.config.BackupDir, backupPath)
			err := MakeDirectories(backupDirPath, info.Mode())
			if err != nil {
				return err
			}

		} else {
			_, err := CopyFile(filepath.Join(siteBuilder.config.BackupDir, backupPath), path)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (siteBuilder *SiteBuilder) clearTargetDir() error {
	return os.RemoveAll(siteBuilder.config.TargetDir)
}

// NEED A BETTER DIRECTORY CRAWLER
func (siteBuilder *SiteBuilder) buildSiteFiles() error {
	siteVars := siteBuilder.siteVars()
	return filepath.Walk(siteBuilder.config.SourceDir, func(path string, info os.FileInfo, _ error) error {
		// check if path should be ignored
		ignore := siteBuilder.ignorePath(path)
		if info.IsDir() {
			if ignore {
				return filepath.SkipDir
			}
			return nil
		}
		if ignore {
			return nil
		}

		// Make a siteFile
		siteFile, err := NewSiteFile(path, siteBuilder, siteVars)
		if err != nil {
			return err
		}

		// Make web file in target directory
		pagePath := siteBuilder.pagePath(path)
		err = MakeDirectories(filepath.Dir(pagePath), os.FileMode(DIRECTORY_PERMISSIONS))
		if err != nil {
			return err
		}
		pageFile, err := os.Create(pagePath)
		if err != nil {
			return err
		}
		defer pageFile.Close()

		_, err = pageFile.WriteString(siteFile.fullContent)
		if err != nil {
			return err
		}

		return nil
	})
}

func (siteBuilder *SiteBuilder) pagePath(path string) string {
	relativePath := strings.TrimPrefix(path, siteBuilder.config.SourceDir)
	return filepath.Join(siteBuilder.config.TargetDir, relativePath)
}

func (siteBuilder *SiteBuilder) siteVars() SiteVars {
	return SiteVars{"url": siteBuilder.config.SiteUrl}
}

func (siteBuilder *SiteBuilder) templatePath(path string) string {
	return filepath.Join(siteBuilder.config.SourceDir, siteBuilder.config.TemplateDir, path)
}

// Janky function...
func (siteBuilder *SiteBuilder) ignorePath(path string) bool {
	relativePath, err := filepath.Rel(siteBuilder.config.SourceDir, path)
	if err != nil {
		panic(err.Error())
	}
	if relativePath == "." {
		return false
	}
	slashRelativePath := strings.TrimPrefix(filepath.ToSlash(relativePath), "/")
	pathParts := strings.Split(slashRelativePath, "/")
	ignorePatterns := []string{"_*", ".*", siteBuilder.config.TemplateDir}
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
