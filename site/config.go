package site

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
)

type siteConfig struct {
	Backup    string
	BackupDir string
	SiteUrl   string
	SourceDir string
	TargetDir string
	LayoutDir string
}

func NewConfig(cmdln map[string]*string) *siteConfig {
	sc := &siteConfig{}
	for key, val := range cmdln {
		ref := reflect.ValueOf(sc).Elem().FieldByName(key)
		ref.SetString(*val)
	}
	return sc.parseConfigFile()
}

func (s *siteConfig) parseConfigFile() *siteConfig {
	sf, err := os.Open("_config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer sf.Close()

	jsonParser := json.NewDecoder(sf)
	err = jsonParser.Decode(s)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
