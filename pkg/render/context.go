package render

import (
	"bytes"
	"github.com/westcoastcode-se/gocms/pkg/content"
	"html/template"
	"io"
	"log"
)

type Context struct {
	Model            *content.Model
	Funcs            template.FuncMap
	templateDatabase TemplateDatabase
}

func (c *Context) RenderView(bw io.Writer, view string, model interface{}) error {
	t := template.New("").Funcs(c.Funcs).Funcs(template.FuncMap{
		"RenderView": func(view string, data interface{}) (template.HTML, error) {
			buf := bytes.NewBuffer([]byte{})
			err := c.RenderView(buf, view, data)

			if err != nil {
				log.Print("failed to render inner template: " + err.Error())
				return "", nil
			}
			return template.HTML(buf.String()), nil
		},
	})

	err := c.templateDatabase.ParseTemplates(t)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(bw, view, model)
	if err != nil {
		log.Print("failed to render inner template: " + err.Error())
		return err
	}
	return nil
}
