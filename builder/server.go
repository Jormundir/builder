package main

import (
	"builder"
	"log"
)

var cmdServer = &Command{
	Name:  "server",
	usage: "[-hi hello] [flags]",
	Short: "short of server command.",
	Long:  `Long of server command.`,
}

func init() {
	cmdServer.Run = runServer
}

func runServer(args []string) {
	config := builder.NewSiteConfig(args)
	site := builder.NewSite(config, nil)
	watcher := site.MakeWatcher(func() {
		site.InitGenerate()
	})
	webserver := site.MakeWebserver()
	if err := watcher.Watch(); err != nil {
		log.Fatalln(err)
	}
	if err := webserver.Serve(); err != nil {
		log.Fatalln(err)
	}
}
