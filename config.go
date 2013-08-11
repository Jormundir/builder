package builder

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

type SiteConfig struct {
	SourceDir string
	TargetDir string
	LayoutDir string

	Url string

	flag flag.FlagSet
}

func NewSiteConfig(args []string) *SiteConfig {
	sc := &SiteConfig{
		SourceDir: ".",
		TargetDir: "_site",
		LayoutDir: "_layouts",
		Url:       "localhost:4000",
	}
	sc.parseConfigFile()
	sc.flag.StringVar(&sc.SourceDir, "source", sc.SourceDir, "Path to project source files.")
	sc.flag.StringVar(&sc.TargetDir, "target", sc.TargetDir, "Directory site will be built in.")
	sc.flag.StringVar(&sc.LayoutDir, "layouts", sc.LayoutDir, "Directory layouts will be looked for in.")
	sc.flag.StringVar(&sc.Url, "url", sc.Url, "The url of your website.")
	sc.flag.Parse(args)
	return sc
}

func (sc *SiteConfig) copy() *SiteConfig {
	config := *sc
	return &config
}

func (sc *SiteConfig) appendPaths(path string) *SiteConfig {
	sc.SourceDir = filepath.Join(sc.SourceDir, path)
	sc.TargetDir = filepath.Join(sc.TargetDir, path)
	sc.LayoutDir = filepath.Join(sc.LayoutDir, path)
	return sc
}

func (sc *SiteConfig) parseConfigFile() {
	sf, err := os.Open("_config.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer sf.Close()

	jsonParser := json.NewDecoder(sf)
	err = jsonParser.Decode(sc)
	if err != nil {
		log.Fatalln(err)
	}
}
