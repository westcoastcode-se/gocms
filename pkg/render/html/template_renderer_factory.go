package html

import (
	"github.com/westcoastcode-se/gocms/pkg/config"
	"github.com/westcoastcode-se/gocms/pkg/content"
	"github.com/westcoastcode-se/gocms/pkg/log"
	"github.com/westcoastcode-se/gocms/pkg/render"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

type TemplateRendererFactory struct {
	ContentRepository content.Repository
	TemplateDatabase  TemplateDatabase
	Config            config.Config
}

func (h *TemplateRendererFactory) NewRenderer(r *http.Request) render.TemplateRenderer {
	uri := r.URL.Path
	user, _ := r.Context().Value(security.SessionKey).(*security.User)
	funcs := template.FuncMap{
		"Navigation": func() *content.Navigation { return &content.Navigation{uri} },
		"Author":     func() bool { return h.Config.Author },
		"Public":     func() bool { return !h.Config.Author },
		"User":       func() *security.User { return user },
		"Search": func(contentType string) []*content.SearchResult {
			return h.ContentRepository.Search(contentType)
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
				log.Warnf(r.Context(), "Could not include file: %s. Reason: %e", path, err)
				return ""
			}
			return template.JS(b)
		},
		"IncludeCSS": func(path string) template.CSS {
			absolutePath := "content" + path
			bytes, err := ioutil.ReadFile(absolutePath)
			if err != nil {
				log.Warnf(r.Context(), "Could not include file: %s. Reason: %e", path, err)
				return ""
			}
			return template.CSS(bytes)
		},
		"IncludeScript": func(path string) template.JS {
			absolutePath := "content" + path
			bytes, err := ioutil.ReadFile(absolutePath)
			if err != nil {
				log.Warnf(r.Context(), "Could not include file: %s. Reason: %e", path, err)
				return ""
			}
			html := template.JS(bytes)
			return html
		},
		"Lookup": func(id string) *content.SearchResult {
			return h.ContentRepository.Lookup(id)
		},
	}

	return &TemplateRenderer{r.Context(), h.TemplateDatabase, funcs}
}
