package render

import (
	"bytes"
	"github.com/westcoastcode-se/gocms/content"
	"html/template"
	"log"
)

type Context struct {
	Model            *content.Model
	Funcs            template.FuncMap
	templateDatabase *TemplateDatabase
}

func (c *Context) GetTemplate() (*template.Template, error) {
	t := template.New("").Funcs(c.Funcs).Funcs(template.FuncMap{
		"RenderView": func(view string, data interface{}) (template.HTML, error) {
			innerTemplate, err := c.GetTemplate()

			if err != nil {
				log.Print("failed to render inner template: " + err.Error())
				return "", nil
			}

			buf := bytes.NewBuffer([]byte{})
			err = innerTemplate.ExecuteTemplate(buf, view, data)
			if err != nil {
				log.Print("failed to render inner template: " + err.Error())
				return "", nil
			}
			return template.HTML(buf.String()), nil
		},
	})

	err := c.templateDatabase.ParseTemplates(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
