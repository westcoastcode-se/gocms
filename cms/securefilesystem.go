package cms

import (
	"errors"
	"net/http"
)

// Customized FileSystem to prevent listing of files for static files
type SecureFileSystem struct {
	fs http.FileSystem
}

func (fs SecureFileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		return nil, errors.New("cannot display directories")
	}

	return f, nil
}

func NewSecureFileSystem(path string) http.FileSystem {
	return &SecureFileSystem{http.Dir(path)}
}
