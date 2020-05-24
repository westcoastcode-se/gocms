package render

import (
	"github.com/westcoastcode-se/gocms/pkg/config"
	"github.com/westcoastcode-se/gocms/pkg/content"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

// Factory used when creating a context. The context contains the neccesary data and functions
// user when rendering a html template
type ContextFactory interface {
	// Create a context used when rendering a html template
	Create(r *http.Request, model *content.Model) *Context
}

// Default context factory
type DefaultContextFactory struct {
	ContentRepository content.Repository
	TemplateDatabase  TemplateDatabase
	Config            config.Config
}

func (d *DefaultContextFactory) Create(r *http.Request, model *content.Model) *Context {
	uri := r.URL.Path
	user, _ := r.Context().Value(security.SessionKey).(*security.User)
	return &Context{
		Model: model,
		Funcs: template.FuncMap{
			"Model":      func() *content.Model { return model },
			"View":       func() string { return model.View },
			"Navigation": func() *content.Navigation { return &content.Navigation{uri} },
			"Author":     func() bool { return d.Config.Author },
			"Public":     func() bool { return !d.Config.Author },
			"User":       func() *security.User { return user },
			"Search": func(contentType string) []*content.SearchResult {
				return d.ContentRepository.Search(contentType)
			},
			"Sort": func(by string, asc string, content []*content.SearchResult) []*content.SearchResult {
				if by == "CreatedAt" {
					if asc == "asc" {
						sort.SliceStable(content, func(i, j int) bool {
							return content[i].Model.CreatedAt.Before(content[j].Model.CreatedAt)
						})
					} else {
						sort.SliceStable(content, func(i, j int) bool {
							return content[j].Model.CreatedAt.Before(content[i].Model.CreatedAt)
						})
					}
				}
				return content
			},
			"Limit": func(limit int, a []*content.SearchResult) []*content.SearchResult {
				return a[0:limit]
			},
			"RenderScript": func(view string) template.JS {
				view = view[:len(view)-5]
				path := "/assets/js/" + view + ".js"
				absolutePath := "./content" + path
				_, err := os.Stat(absolutePath)
				if os.IsNotExist(err) {
					return ""
				}

				b, err := ioutil.ReadFile(absolutePath)
				if err != nil {
					log.Printf("Could not include file: %s. Reason: %e", path, err)
					return ""
				}
				return template.JS(b)
			},
			"IncludeCSS": func(path string) template.CSS {
				absolutePath := "content" + path
				bytes, err := ioutil.ReadFile(absolutePath)
				if err != nil {
					log.Printf("Could not include file: %s. Reason: %e", path, err)
					return ""
				}
				return template.CSS(bytes)
			},
			"IncludeScript": func(path string) template.JS {
				absolutePath := "content" + path
				bytes, err := ioutil.ReadFile(absolutePath)
				if err != nil {
					log.Printf("Could not include file: %s. Reason: %e", path, err)
					return ""
				}
				html := template.JS(bytes)
				return html
			},
			"Lookup": func(id string) *content.SearchResult {
				return d.ContentRepository.Lookup(id)
			},
		},
		templateDatabase: d.TemplateDatabase,
	}
}
