package site

import (
	"github.com/howeyc/fsnotify"
	"log"
)

type watcher struct {
	site *Site
}

func newWatcher(s *Site) *watcher {
	return &watcher{s}
}

func (w *watcher) watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println(ev.Name + " detected, rebuilding site..")
				w.site.buildMode()
				w.site.Init()
				w.site.Build()
				w.site.webMode()
				w.site.Init()
			case er := <-watcher.Error:
				log.Println(er)
				break
			}
		}
	}()

	for _, dir := range w.site.dirs {
		err := watcher.Watch(dir)
		if err != nil {
			return err
		}
	}

	return nil
}
