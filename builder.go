package main

import (
	"builder/site"
	"flag"
)

var options = map[string]*string{
	"Backup":    flag.String("Backup", "true", "Indicate if you want target directory to be backed up to a different folder before being wiped."),
	"BackupDir": flag.String("BackupDir", "_backup", "Directory to backup target directory to."),
	"SiteUrl":   flag.String("SiteUrl", "http://localhost:1337", "Base url of your site."),
	"SourceDir": flag.String("SourceDir", ".", "Path to your project's source directory."),
	"TargetDir": flag.String("TargetDir", "_site", "Path to where your site will be built."),
	"LayoutDir": flag.String("LayoutDir", "_layouts", "Path to layouts directory."),
}

func main() {
	flag.Parse()
	s := site.NewSite(options)
	s.Execute()
}
