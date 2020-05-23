package immediate

import (
	"github.com/westcoastcode-se/gocms/render"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Database struct {
	rootPath string
}

func (f *Database) ParseTemplates(original *template.Template) error {
	pfx := len(f.rootPath) + 1
	templates := make(map[string]string)
	err := filepath.Walk(f.rootPath, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			name = filepath.ToSlash(name)
			templates[name] = string(b)
			log.Printf("Loaded template file: %s", name)
		}
		return nil
	})

	if err != nil {
		return err
	}

	for key, value := range templates {
		t := original.New(key)
		_, err := t.Parse(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Database) FindTemplate(path string) (string, error) {
	_, err := os.Stat(f.rootPath + path)
	if err != nil {
		return "", &render.TemplateNotFound{Path: path}
	}

	bytes, err := ioutil.ReadFile(f.rootPath + path)
	if err != nil {
		return "", &render.TemplateNotFound{Path: path}
	}

	return string(bytes), nil
}

// Create a database based on the template directory instead of keeping it in memory
func NewFileSystemTemplateDatabase(rootPath string) render.TemplateDatabase {
	impl := &Database{rootPath: rootPath}
	return impl
}
