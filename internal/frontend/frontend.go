package frontend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// GetJSConfigHandler returns a HandlerFunc to deal with the JS const values.
func GetJSConfigHandler(assetsPath string, varMapping map[string]string) (http.HandlerFunc, error) {
	f, err := os.Open(assetsPath + "/vars.js")
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	for oldStr, newStr := range varMapping {
		b = bytes.ReplaceAll(b, []byte(oldStr), []byte(newStr))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}, nil
}

// GetFileSystemHandler get the file system http handler
func GetFileSystemHandler(directory, path string) http.HandlerFunc {
	fileServer := http.FileServer(FileSystem{http.Dir(directory)})
	return http.StripPrefix(strings.TrimRight(path, "/"), fileServer).ServeHTTP
}

// FileSystem custom file system handler
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
