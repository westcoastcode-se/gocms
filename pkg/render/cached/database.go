package cached

import (
	"context"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"github.com/westcoastcode-se/gocms/pkg/log"
	"github.com/westcoastcode-se/gocms/pkg/render"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Database struct {
	rootPath  string
	mux       sync.Mutex
	Templates map[string]string
}

func (f *Database) ParseTemplates(original *template.Template) error {
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
func (f *Database) FindTemplate(path string) (string, error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	if result, ok := f.Templates[path]; ok {
		return result, nil
	}
	return "", &render.TemplateNotFound{Path: path}
}

func (f *Database) load(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Infof("Loading template files from %s", f.rootPath)
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
			log.Infof("Loaded template file: %s", name)
		}
		return nil
	})
	f.mux.Lock()
	defer f.mux.Unlock()
	f.Templates = templates
	return err
}

func (f *Database) OnEvent(ctx context.Context, e interface{}) error {
	if _, ok := e.(*event.Checkout); ok {
		if err := f.load(ctx); err != nil {
			return err
		}
	}
	return nil
}

func NewDatabase(bus *event.Bus, rootPath string) render.TemplateDatabase {
	impl := &Database{
		rootPath:  rootPath,
		mux:       sync.Mutex{},
		Templates: make(map[string]string),
	}
	if len(rootPath) > 0 {
		err := impl.load(context.Background())
		if err != nil {
			panic(err)
		}
	}
	bus.AddListener(impl)
	return impl
}
