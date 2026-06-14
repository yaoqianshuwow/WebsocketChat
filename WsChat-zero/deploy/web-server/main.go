package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	addr = flag.String("addr", ":80", "listen address")
	dir  = flag.String("dir", "./web-dist", "static file directory")
)

func main() {
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("resolve static directory failed: %v", err)
	}

	fileServer := http.FileServer(http.Dir(absDir))
	indexPath := filepath.Join(absDir, "index.html")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := filepath.Clean("/" + r.URL.Path)
		if cleanPath == "/api" || strings.HasPrefix(cleanPath, "/api/") || cleanPath == "/wss" {
			http.NotFound(w, r)
			return
		}

		target := filepath.Join(absDir, filepath.FromSlash(strings.TrimPrefix(cleanPath, "/")))
		if info, err := os.Stat(target); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, indexPath)
	})

	log.Printf("static server listening on %s, serving %s", *addr, absDir)
	log.Fatal(http.ListenAndServe(*addr, handler))
}
