package builder

import (
	"github.com/howeyc/fsnotify"
	"log"
)

type watcher struct {
	site    *Site
	onEvent func()
}

func newWatcher(s *Site, onEvent func()) *watcher {
	return &watcher{
		site:    s,
		onEvent: onEvent,
	}
}

func (w *watcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println(ev.Name + " detected.")
				w.onEvent()
			case er := <-watcher.Error:
				log.Println(er)
				break
			}
		}
	}()

	for _, dir := range w.site.getDirs() {
		err := watcher.Watch(dir)
		if err != nil {
			return err
		}
	}

	return nil
}
