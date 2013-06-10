package site

import (
	"encoding/json"
	"flag"
	"os"
)

type siteConfig struct {
	Backup    bool
	BackupDir string
	SiteUrl   string
	SourceDir string
	TargetDir string
	LayoutDir string
}

func newSiteConfig() (*siteConfig, error) {
	// Make config with defaults
	siteConfig := &siteConfig{
		Backup:    true,
		BackupDir: "_backup",
		SiteUrl:   "http://localhost:4000",
		SourceDir: ".",
		TargetDir: "_site",
		LayoutDir: "layouts",
	}
	// Override defaults with config file declarations
	err := siteConfig.parseConfigFile()
	if err != nil {
		return nil, err
	}
	// Override everything with command line arguments
	siteConfig.parseCommandLine()
	return siteConfig, nil
}

func (siteConfig *siteConfig) parseConfigFile() error {
	configFile, err := os.Open("_config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(siteConfig)
	if err != nil {
		return err
	}
	return nil
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
