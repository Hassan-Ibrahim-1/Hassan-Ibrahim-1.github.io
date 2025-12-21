package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	allowed := map[string]bool{
		"index.html":                true,
		"theme.css":                 true,
		"hassan-ibrahim-resume.pdf": true,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip leading slash
		name := filepath.Base(r.URL.Path)

		// Reject directories or empty paths
		if name == "." || name == "/" {
			http.NotFound(w, r)
			return
		}

		// Enforce whitelist
		if !allowed[name] {
			http.NotFound(w, r)
			return
		}

		// Serve from current directory
		f, err := os.Open(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		http.ServeContent(w, r, name, getModTime(f), f)
	})

	log.Println("Serving on http://localhost:5100")
	log.Fatal(http.ListenAndServe(":5100", handler))
}

func getModTime(f *os.File) (t time.Time) {
	if info, err := f.Stat(); err == nil {
		t = info.ModTime()
	}
	return
}
