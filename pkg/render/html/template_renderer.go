package html

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/westcoastcode-se/gocms/pkg/log"
	"html/template"
	"io"
)

type TemplateRenderer struct {
	// The underlying context responsible for creating this renderer
	ctx context.Context
	// Database containing all templates
	templateDatabase TemplateDatabase
	//
	funcs template.FuncMap
}

func (h *TemplateRenderer) RenderView(writer io.Writer, view string, model interface{}) error {
	t := template.New("")
	t = t.Funcs(h.funcs).Funcs(template.FuncMap{
		"Model": func() interface{} { return model },
		"RenderView": func(childView string, childModel interface{}) (template.HTML, error) {
			buf := bytes.NewBuffer([]byte{})
			err := h.RenderView(buf, childView, childModel)
			if err != nil {
				log.Errorf(h.ctx, "Failed to render child view %s: %e", childView, err)
				return "", err
			}
			return template.HTML(buf.String()), nil
		},
	})

	err := h.templateDatabase.ParseTemplates(t)
	if err != nil {
		return errors.New(fmt.Sprintf("could not parse templates: %e", err))
	}

	err = t.ExecuteTemplate(writer, view, model)
	if err != nil {
		return errors.New(fmt.Sprintf("could not execute template with template "+view+": %e", err))
	}

	return nil
}
