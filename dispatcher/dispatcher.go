package dispatcher

import (
	"builder/site"
	"fmt"
	"github.com/howeyc/fsnotify"
)

type Dispatcher struct {
	site   *site.Site
	config *site.SiteConfig
}

func NewDispatcher() (*Dispatcher, error) {
	config, err := site.NewSiteConfig()
	if err != nil {
		return nil, err
	}
	site, err := site.NewSite(config)
	if err != nil {
		return nil, err
	}
	return &Dispatcher{config: config, site: site}, nil
}

func (dispatcher *Dispatcher) Dispatch() error {
	fmt.Print("building site...")
	errs := dispatcher.site.BuildSite()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}
	fmt.Println(" done")
	fmt.Print("watching source directory...")
	dispatcher.watch()
	fmt.Println(" done")
	fmt.Print("starting webserver...")
	server := dispatcher.site.MakeWebServer()
	err := server.Serve()
	if err != nil {
		return err
	}
	return nil
}

func (dispatcher *Dispatcher) watch() error {
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return
		}

		done := make(chan bool)

		go func() {
			for {
				select {
				case ev := <-watcher.Event:
					fmt.Println(ev, " detected. Rebuilding site and restarting server.")
					dispatcher.site, err = site.NewSite(dispatcher.config)
					if err != nil {
						fmt.Println(err.Error())
						<-done
					}
					dispatcher.Dispatch()
				case err := <-watcher.Error:
					fmt.Println(err.Error())
					<-done
				}
			}
		}()

		for _, dir := range dispatcher.site.Dirs() {
			err := watcher.Watch(dir)
			if err != nil {
				fmt.Println(err.Error())
				<-done
			}
		}

		<-done
		watcher.Close()
		return
	}()
	return nil
}
