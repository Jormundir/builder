package site

import (
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type webserver struct {
	site *Site
}

func newWebserver(s *Site) *webserver {
	return &webserver{s}
}

func (s webserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.FromSlash(strings.TrimLeft(r.URL.Path, "/"))
	rp, ok := s.site.pages[path]
	if !ok {
		log.Println(path + " requested but does not exist.")
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(rp.htmlext))
	io.WriteString(w, rp.page)
}

func (s *webserver) serve() error {
	err := http.ListenAndServe("localhost:1337", s)
	if err != nil {
		return err
	}
	return nil
}
