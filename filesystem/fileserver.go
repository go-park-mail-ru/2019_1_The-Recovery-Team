package filesystem

import (
	"net/http"
)

// InitFileserverHandler initialize handler for serving files
func InitFileserverHandler(path, dir string) http.HandlerFunc {
	fileserver := func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/upload/", http.FileServer(http.Dir("./upload/"))).ServeHTTP(w, r)
	}
	return fileserver
}
