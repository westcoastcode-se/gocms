package render

import "html/template"

// Database used for managing page templates
type TemplateDatabase interface {
	// Parse all templates and put the content into the supplied template context
	ParseTemplates(original *template.Template) error

	// Search for a template at the supplied path
	FindTemplate(path string) (string, error)
}
