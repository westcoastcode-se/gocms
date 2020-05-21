package render

import (
	"github.com/westcoastcode-se/gocms/event"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type GitControlledTemplateDatabase struct {
	rootPath  string
	mux       sync.Mutex
	Templates map[string]string
}

func (f *GitControlledTemplateDatabase) ParseTemplates(original *template.Template) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for key, value := range f.Templates {
		t := original.New(key)
		_, err := t.Parse(value)
		if err != nil {
			return err
		}
	}
	return nil
}

// Fetch a template based on it's path
func (f *GitControlledTemplateDatabase) FindTemplate(path string) (string, error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	if result, ok := f.Templates[path]; ok {
		return result, nil
	}
	return "", &TemplateNotFound{path}
}

func (f *GitControlledTemplateDatabase) load() error {
	log.Printf("Loading template files from %s", f.rootPath)

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
	f.mux.Lock()
	defer f.mux.Unlock()
	f.Templates = templates
	return err
}

func (f *GitControlledTemplateDatabase) OnEvent(e interface{}) error {
	if _, ok := e.(*event.Checkout); ok {
		if err := f.load(); err != nil {
			return err
		}
	}
	return nil
}

func NewTemplateDatabase(bus *event.Bus, rootPath string) TemplateDatabase {
	impl := &GitControlledTemplateDatabase{
		rootPath:  rootPath,
		mux:       sync.Mutex{},
		Templates: make(map[string]string),
	}
	if len(rootPath) > 0 {
		err := impl.load()
		if err != nil {
			panic(err)
		}
	}
	bus.AddListener(impl)
	return impl
}
