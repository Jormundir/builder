package site

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WebServer struct {
	pages map[string]*page
	port  string
}

func (server *WebServer) lookupContentType(ext string) string {
	switch ext {
	case ".css":
		return "text/css"
	case ".js":
		return "application/x-javascript"
	default:
		return "text/html"
	}
}

func (server WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	requestedPage, ok := server.pages[path]
	if !ok {
		fmt.Println(path + " does not exist.")
		return
	}
	w.Header().Set("Content-Type", server.lookupContentType(requestedPage.ext))
	io.WriteString(w, requestedPage.fullHtmlContent)
}

func (server *WebServer) Serve() error {
	err := http.ListenAndServe("localhost:"+server.port, server)
	if err != nil {
		return err
	}
	return nil
}
