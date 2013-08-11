package main

import (
	"builder"
)

var cmdBuild = &Command{
	Name:  "build",
	usage: "[-hi hello] [flags]",
	Short: "short of build command.",
	Long:  `Long of build command.`,
}

func init() {
	cmdBuild.Run = runBuild
}

func runBuild(args []string) {
	config := builder.NewSiteConfig(args)
	site := builder.NewSite(config, nil)
	site.GenerateFiles()
}
