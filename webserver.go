package builder

import (
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
)

type webserver struct {
	site *Site
}

func newWebServer(s *Site) *webserver {
	return &webserver{s}
}

func (ws *webserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := ws.site.webpage(strings.TrimLeft(r.URL.Path, "/"))
	if page == nil {
		log.Println(r.URL.Path + " requested but does not exist.")
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(page.webExt))
	io.WriteString(w, page.content)
}

func (ws *webserver) Serve() error {
	err := http.ListenAndServe(ws.site.config.Url, ws)
	if err != nil {
		return err
	}
	return nil
}
